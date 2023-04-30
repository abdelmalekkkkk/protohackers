package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/abdelmalekkkkk/protohackers/server"
	"github.com/wangjia184/sortedset"
)

type MessageType byte

const (
	Insert MessageType = 'I'
	Query  MessageType = 'Q'
)

type Message struct {
	Type        MessageType
	FirstValue  int32
	SecondValue int32
}

type Entry struct {
	Price     int32
	Timestamp int32
}

type PricesController struct {
	LastID uint64
	Prices *sortedset.SortedSet
}

func main() {
	port := flag.String("port", "69", "Which port to listen to")
	host := flag.String("host", "0.0.0.0", "Which host to bind the listener to")

	flag.Parse()

	server, err := server.NewServer(server.ServerConfig{
		Host:           *host,
		Port:           *port,
		RequestHandler: assetRequestHandler,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func assetRequestHandler(conn net.Conn) {
	defer conn.Close()

	controller := PricesController{
		Prices: sortedset.New(),
	}

	buf := make([]byte, 9)

	for {
		_, err := io.ReadFull(conn, buf)

		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
			} else {
				fmt.Println(fmt.Errorf("could not read data: %w", err))
			}
			break
		}

		fmt.Println(buf)

		message, err := controller.parseAssetRequest(buf)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(message)

		if message.Type == Insert {
			controller.insertPrice(message)
			continue
		}

		average := controller.queryPrice(message)
		fmt.Println(average)
		responseBytes := toByteArray(int32(average))
		conn.Write(responseBytes[:])
	}
}

func (c *PricesController) parseAssetRequest(data []byte) (*Message, error) {
	if len(data) != 9 {
		return nil, fmt.Errorf("invalid bytes size")
	}

	message := new(Message)

	if data[0] == byte(Insert) {
		message.Type = Insert
	} else if data[0] == byte(Query) {
		message.Type = Query
	} else {
		return nil, fmt.Errorf("incorrect message type")
	}

	message.FirstValue = int32(binary.BigEndian.Uint32(data[1:5]))
	message.SecondValue = int32(binary.BigEndian.Uint32(data[5:9]))

	return message, nil
}

func (c *PricesController) insertPrice(message *Message) {
	ID := atomic.AddUint64(&c.LastID, 1)
	key := strconv.Itoa(int(ID))

	c.Prices.AddOrUpdate(key, sortedset.SCORE(message.FirstValue), message.SecondValue)
}

func (c *PricesController) queryPrice(message *Message) float64 {
	min := message.FirstValue
	max := message.SecondValue

	if max < min {
		return 0
	}

	average := 0.00

	prices := c.Prices.GetByScoreRange(sortedset.SCORE(min), sortedset.SCORE(max), nil)
	size := len(prices)

	for _, price := range prices {
		convertedPrice := price.Value.(int32)
		average += float64(convertedPrice) / float64(size)
	}

	return average
}

func toByteArray(i int32) (arr [4]byte) {
	binary.BigEndian.PutUint32(arr[0:4], uint32(i))
	return
}
