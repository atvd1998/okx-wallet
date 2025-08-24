package http

import (
	"net/http"
)

type Config struct {
	HttpAddress string
}

func NewHttpServer(config *Config) *http.Server {
	mux := http.NewServeMux()
	handler := http.Handler(mux)
	return &http.Server{
		Addr:    config.HttpAddress,
		Handler: handler,
	}
}
