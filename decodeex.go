package traefik_middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Config struct {
	xDecodePath string
	xRequestID  string
}

func CreateConfig() *Config {
	return &Config{xDecodePath: "x-decode-path", xRequestID: "x-request-id"}
}

func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Add(config.xDecodePath, r.URL.Path)

		if r.Header.Get(config.xRequestID) == "" {
			requestID := uuid.Must(uuid.NewRandom()).String()
			r.Header.Add(config.xRequestID, requestID)
			rw.Header().Add(config.xRequestID, requestID)
		} else {
			rw.Header().Add(config.xRequestID, r.Header.Get(config.xRequestID))
		}
		next.ServeHTTP(rw, r)
	})
	return handler, nil
}
