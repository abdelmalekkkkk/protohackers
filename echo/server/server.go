package server

import (
	"fmt"
	"io/ioutil"
	"net"
)

type ServerConfig struct {
	Port string
	Host string
}

type server struct {
	Config ServerConfig
}

func NewServer(config ServerConfig) (*server, error) {
	if config.Port == "" {
		return nil, fmt.Errorf("the port cannot be empty")
	}

	if config.Port == "" {
		return nil, fmt.Errorf("the host cannot be empty")
	}

	return &server{
		Config: config,
	}, nil
}

func (server *server) Run() error {
	addr := server.Config.Host + ":" + server.Config.Port

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start listening at %s: %w", addr, err)
	}

	defer listen.Close()
	fmt.Printf("Started listening on %s\n", addr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(fmt.Errorf("could not accept connection: %w", err))
		}

		server.HandleConnection(conn)
	}
}

func (server *server) HandleConnection(conn net.Conn) {
	data, err := ioutil.ReadAll(conn)

	if err != nil {
		fmt.Println(fmt.Errorf("could not read data: %w", err))
	}

	conn.Write(data)
	conn.Close()
}
