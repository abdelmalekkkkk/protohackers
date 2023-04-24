package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/abdelmalekkkkk/protohackers/server"
)

func main() {
	port := flag.String("port", "7", "Which port to listen to")
	host := flag.String("host", "0.0.0.0", "Which host to bind the listener to")

	flag.Parse()

	server, err := server.NewServer(server.ServerConfig{
		Host: *host,
		Port: *port,
		RequestHandler: func(conn net.Conn) {
			data, err := ioutil.ReadAll(conn)

			if err != nil {
				fmt.Println(fmt.Errorf("could not read data: %w", err))
			}

			conn.Write(data)
			conn.Close()
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
