package traefik_middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	decodeex "github.com/decodeex/traefik-plugin"
)

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}

func TestXDecodePath(t *testing.T) {
	cfg := decodeex.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := decodeex.New(ctx, next, cfg, "decodeex-plugin")
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/v1/demo/users", nil)
	handler.ServeHTTP(recorder, req)
	assertHeader(t, req, "x-decode-path", "/v1/demo/users")
}
