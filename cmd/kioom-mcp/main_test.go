package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	token := "my-secret-token"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	authed := authMiddleware(token, handler)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
	}{
		{
			name: "valid bearer token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/", nil)
				req.Header.Set("Authorization", "Bearer my-secret-token")
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid bearer token case insensitive prefix",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/", nil)
				req.Header.Set("Authorization", "bearer my-secret-token")
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid bearer token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/", nil)
				req.Header.Set("Authorization", "Bearer wrong-token")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "missing authorization",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "valid token via query param token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/?token=my-secret-token", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid token via query param auth",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/?auth=my-secret-token", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid token via query param",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "http://example.com/?token=wrong-token", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()
			rec := httptest.NewRecorder()
			authed.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("Accepted"))
	})

	logged := loggingMiddleware(handler)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	rec := httptest.NewRecorder()
	logged.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
}

type mockFlusherResponseWriter struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (m *mockFlusherResponseWriter) Flush() {
	m.flushed = true
}

func TestLoggingResponseWriterFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	w := &mockFlusherResponseWriter{ResponseRecorder: rec}
	lrw := &loggingResponseWriter{ResponseWriter: w}
	lrw.Flush()

	if !w.flushed {
		t.Error("expected Flush to be called on underlying ResponseWriter")
	}
}
