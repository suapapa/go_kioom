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
	"log"
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
	log.SetPrefix("kioom-mcp: ")
	log.SetFlags(0)

	// Load environment variables from .env if present
	if err := kioomenv.LoadEnvFile(".env"); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	transport := flag.String("transport", "stdio", "Transport: stdio or sse")
	listen := flag.String("listen", "127.0.0.1:8765", "Listen address for -transport=sse")
	ssePath := flag.String("sse-path", "/", "HTTP path prefix for -transport=sse (non-root paths get a trailing / for POST session routing)")
	sseToken := flag.String("sse-token", os.Getenv("KIOOM_MCP_SSE_TOKEN"), "Bearer token to authenticate incoming SSE client connections")

	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var opts []kioom.Option
	if cfg.Mock {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)
	if cfg.Token != "" {
		client.SetToken(cfg.Token)
	}

	srv := mcpkioom.NewServer(client, nil, nil)

	switch strings.ToLower(strings.TrimSpace(*transport)) {
	case "stdio":
		rootCtx := signalContext(context.Background())
		defer rootCtx.stop()
		if err := srv.Run(rootCtx.ctx, &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}
	case "sse":
		if err := runSSE(srv, *listen, normalizeSSEPath(*ssePath), *sseToken); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown -transport %q: use stdio or sse", *transport)
	}
}

func loadConfig() (kioomenv.Config, error) {
	c := kioomenv.Load(os.Getenv)
	if err := c.RequireAppKeys(); err != nil {
		return kioomenv.Config{}, err
	}
	return c, nil
}

func runSSE(mcpsrv *mcp.Server, listenAddr, path, token string) error {
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
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	rootCtx := signalContext(context.Background())
	defer rootCtx.stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on http://%s%s (SSE transport); shutdown with SIGINT/SIGTERM", addr, path)
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
