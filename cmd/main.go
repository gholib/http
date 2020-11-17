package main

import (
	"net"
	"net/http"
	"os"

	"github.com/gholib/http/cmd/app"
	"github.com/gholib/http/pkg/banners"
)

func main() {
	//обьявляем порт и хост
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(h, p string) error {
	mux := http.NewServeMux()

	bannerSvc := banners.NewService()

	server := app.NewServer(mux, bannerSvc)

	server.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(h, p),
		Handler: server,
	}
	return srv.ListenAndServe()
}
