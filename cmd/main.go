package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gholib/http/pkg/server"
)

func main() {
	host := "0.0.0.0"
	port := "8080"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) error {
	srv := server.NewServer(net.JoinHostPort(host, port))
	body := "hello"
	srv.Register("/api/category{category}/{id}", func(req *server.Request) {
		_, err := req.Conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" + body,
		))
		if err != nil {
			log.Print(err)
		}
	})
	return srv.Start()
}
