package traefik_middleware

import (
	"context"
	"net/http"
)

type Config struct {
	xDecodePath string
}

func CreateConfig() *Config {
	return &Config{xDecodePath: "x-decode-path"}
}

func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		xpath := r.URL.Path
		r.Header.Add(config.xDecodePath, xpath)
		next.ServeHTTP(rw, r)
	})
	return handler, nil
}
