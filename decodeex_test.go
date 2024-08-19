package traefik_middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	decodeex "github.com/decodeex/traefik_middleware"
	"github.com/google/uuid"
)

func assertResponseHeader(t *testing.T, resp *http.Response, key, expected string) {
	t.Helper()

	if expected != "" && resp.Header.Get(key) != expected {
		t.Errorf("invalid header value expect: %s; got %s", expected, resp.Header.Get(key))
	}
}
func assertReqHeader(t *testing.T, req *http.Request, key, expected string) string {
	t.Helper()

	if expected != "" && req.Header.Get(key) != expected {
		t.Errorf("invalid header value expect:%s; got %s", expected, req.Header.Get(key))
	}
	return req.Header.Get(key)
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
	assertReqHeader(t, req, "x-decode-path", "/v1/demo/users")

	// 请求头里不预先携带 x-request-id 头
	recorder2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/v1/demo/users", nil)
	handler.ServeHTTP(recorder2, req2)
	xRequestID := assertReqHeader(t, req2, "x-request-id", "")
	assertResponseHeader(t, recorder2.Result(), "x-request-id", xRequestID)

	// 请求头里预先携带 x-request-id 头
	xRequestID3 := uuid.Must(uuid.NewRandom()).String()
	recorder3 := httptest.NewRecorder()
	req3, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/v1/demo/users", nil)
	req3.Header.Add("x-request-id", xRequestID3)
	handler.ServeHTTP(recorder3, req3)
	assertReqHeader(t, req3, "x-request-id", xRequestID3)
	assertResponseHeader(t, recorder3.Result(), "x-request-id", xRequestID3)
}
