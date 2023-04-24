package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"

	"github.com/abdelmalekkkkk/protohackers/server"
)

type PrimeRequest struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type PrimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	port := flag.String("port", "7", "Which port to listen to")
	host := flag.String("host", "0.0.0.0", "Which host to bind the listener to")

	flag.Parse()

	server, err := server.NewServer(server.ServerConfig{
		Host:           *host,
		Port:           *port,
		RequestHandler: primeRequestHandler,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func primeRequestHandler(conn net.Conn) {
	defer conn.Close()

	buf := bufio.NewReader(conn)

	for {
		bytes, err := buf.ReadBytes('\n')

		if err != nil {
			fmt.Println(fmt.Errorf("could not read data: %w", err))
			break
		}

		request, valid := verifyRequest(bytes)

		if !valid {
			// Send malformed response (a PI character)
			conn.Write([]byte{0xe3})
			break
		}

		response, err := generateResponse(request)

		if err != nil {
			fmt.Println(err)
			continue
		}

		conn.Write(response)
	}
}

func generateResponse(request *PrimeRequest) ([]byte, error) {
	response := PrimeResponse{
		Method: "isPrime",
	}

	if isPrime(*request.Number) {
		response.Prime = true
	}

	bytes, err := json.Marshal(response)

	if err != nil {
		return nil, err
	}

	return append(bytes, '\n'), nil
}

func verifyRequest(data []byte) (*PrimeRequest, bool) {
	var primeRequest PrimeRequest

	err := json.Unmarshal(data, &primeRequest)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if primeRequest.Method == nil || primeRequest.Number == nil || *primeRequest.Method != "isPrime" {
		return nil, false
	}

	return &primeRequest, true
}

func isPrime(number float64) bool {
	if math.Trunc(number) != number {
		return false
	}

	return big.NewInt(int64(int(number))).ProbablyPrime(0)
}
