package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func main() {
	server, err := net.ResolveTCPAddr("tcp", "localhost:69")

	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, server)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	go func() {
		for {
			buf := make([]byte, 4)

			_, err := conn.Read(buf)

			if err != nil {
				fmt.Println(fmt.Errorf("could not read data: %w", err))
				continue
			}

			fmt.Printf("\nreceived %v\n", buf)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your request in the following format TYPE,NUMBER,NUMBER: ")
		text, _ := reader.ReadString('\n')
		parts := strings.Split(strings.TrimSpace(text), ",")
		rType := parts[0]
		rInt1, _ := strconv.Atoi(parts[1])
		rInt2, _ := strconv.Atoi(parts[2])

		request := make([]byte, 0, 9)
		request = append(request, []byte(rType)[0])
		request = append(request, IntToByteArray(int32(rInt1))...)
		request = append(request, IntToByteArray(int32(rInt2))...)

		for _, b := range request {
			conn.Write([]byte{b})
		}
	}

}

func IntToByteArray(num int32) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[size-i-1] = byt
	}
	return arr
}
