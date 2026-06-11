// Command kioom-mcp runs a Model Context Protocol server over the Kiwoom REST client.
//
// Transports:
//   - stdio (default): for Cursor/Claude Desktop style local integration.
//   - sse: HTTP Server-Sent Events per MCP spec (GET stream + POST messages).
//
// Credentials use the same environment variables as kioom-cli: KIOOM_APP_KEY,
// KIOOM_SECRET_KEY, optional KIOOM_TOKEN, optional KIOOM_MOCK=true.
package main

import (
	"context"
	"crypto/subtle"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
	"github.com/suapapa/go_kioom/internal/mcpkioom"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport: stdio or sse")
	listen := flag.String("listen", "127.0.0.1:8765", "Listen address for -transport=sse")
	ssePath := flag.String("sse-path", "/", "HTTP path prefix for -transport=sse (non-root paths get a trailing / for POST session routing)")
	sseToken := flag.String("sse-token", os.Getenv("KIOOM_MCP_SSE_TOKEN"), "Bearer token to authenticate incoming SSE client connections")
	logFormat := flag.String("log-format", "text", "Log output format: text (human-readable) or json")

	flag.Parse()

	logger, err := setupLogger(*logFormat, os.Stderr)
	if err != nil {
		slog.Error("invalid logging configuration", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)

	// Load environment variables from .env if present
	if err := kioomenv.LoadEnvFile(".env"); err != nil {
		logger.Warn("failed to load .env file", "error", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		fatal(logger, "configuration error", "error", err)
	}

	var opts []kioom.Option
	if cfg.Mock {
		opts = append(opts, kioom.WithMockDomain())
	}
	opts = append(opts, kioom.WithAutoToken())
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)
	if cfg.Token != "" {
		client.SetToken(cfg.Token)
	}

	srv := mcpkioom.NewServer(client, nil, nil)

	sessionIDs := &sessionLogRegistry{}

	srv.AddReceivingMiddleware(func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			sess := req.GetSession()
			sessionID := sessionIDs.resolve(sess)
			if method == "initialize" && sessionID != "" {
				logger.Info("MCP session started", "session", sessionID)
			}

			var toolName string
			if method == "tools/call" {
				if rawParams, ok := req.GetParams().(*mcp.CallToolParamsRaw); ok {
					toolName = rawParams.Name
				}
			}

			if toolName != "" {
				logger.Info("MCP request", "tool", toolName, "session", sessionID)
			} else {
				logger.Info("MCP request", "method", method, "session", sessionID)
			}

			res, err := next(ctx, method, req)
			if err != nil {
				if toolName != "" {
					logger.Error("MCP response failed", "tool", toolName, "session", sessionID, "error", err)
				} else {
					logger.Error("MCP response failed", "method", method, "session", sessionID, "error", err)
				}
			} else {
				status := "success"
				if res != nil {
					if toolRes, ok := res.(*mcp.CallToolResult); ok && toolRes.IsError {
						status = "tool_error"
					}
				}
				if toolName != "" {
					logger.Info("MCP response", "tool", toolName, "session", sessionID, "status", status)
				} else {
					logger.Info("MCP response", "method", method, "session", sessionID, "status", status)
				}
			}
			return res, err
		}
	})

	switch strings.ToLower(strings.TrimSpace(*transport)) {
	case "stdio":
		rootCtx := signalContext(context.Background())
		defer rootCtx.stop()
		if err := srv.Run(rootCtx.ctx, &mcp.StdioTransport{}); err != nil {
			fatal(logger, "stdio transport stopped", "error", err)
		}
	case "sse":
		if err := runSSE(logger, srv, *listen, normalizeSSEPath(*ssePath), *sseToken); err != nil {
			fatal(logger, "sse transport stopped", "error", err)
		}
	default:
		fatal(logger, "unknown transport", "transport", *transport)
	}
}

func loadConfig() (kioomenv.Config, error) {
	c := kioomenv.Load(os.Getenv)
	if err := c.RequireAppKeys(); err != nil {
		return kioomenv.Config{}, err
	}
	return c, nil
}

func runSSE(logger *slog.Logger, mcpsrv *mcp.Server, listenAddr, path, token string) error {
	h := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server { return mcpsrv }, nil)
	mux := http.NewServeMux()
	if token != "" {
		mux.Handle(path, authMiddleware(token, h))
	} else {
		mux.Handle(path, h)
	}

	addr := strings.TrimSpace(listenAddr)
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	rootCtx := signalContext(context.Background())
	defer rootCtx.stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", addr, "path", path, "transport", "sse")
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-rootCtx.ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		if err := <-errCh; err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

func authMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		reqToken := ""
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			reqToken = authHeader[7:]
		}

		if reqToken == "" {
			reqToken = r.URL.Query().Get("token")
		}
		if reqToken == "" {
			reqToken = r.URL.Query().Get("auth")
		}

		if subtle.ConstantTimeCompare([]byte(reqToken), []byte(token)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type signalStop struct {
	ctx  context.Context
	stop func()
}

func signalContext(parent context.Context) signalStop {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()
	return signalStop{
		ctx: ctx,
		stop: func() {
			signal.Stop(ch)
			cancel()
		},
	}
}

func normalizeSSEPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	// Go 1.22+ ServeMux matches subtrees with a trailing slash; SSE POSTs use ?sessionid= on the same path.
	if p != "/" && !strings.HasSuffix(p, "/") {
		p += "/"
	}
	return p
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Flush() {
	if fl, ok := lrw.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		isSSE := r.Method == "GET" && strings.Contains(r.Header.Get("Accept"), "text/event-stream")

		if isSSE {
			logger.Info("SSE client connecting", "remote", r.RemoteAddr, "uri", r.RequestURI)
		} else {
			logger.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
		}

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		if isSSE {
			logger.Info("SSE client disconnected", "remote", r.RemoteAddr, "uri", r.RequestURI, "duration", duration)
		} else {
			logger.Info("HTTP response", "status", lrw.statusCode, "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr, "duration", duration)
		}
	})
}
