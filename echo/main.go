package main

import (
	"flag"
	"log"

	"github.com/abdelmalekkkkk/protohackers/echo/server"
)

func main() {
	port := flag.String("port", "7", "Which port to listen to")
	host := flag.String("host", "0.0.0.0", "Which host to bind the listener to")

	flag.Parse()

	server, err := server.NewServer(server.ServerConfig{
		Host: *host,
		Port: *port,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
