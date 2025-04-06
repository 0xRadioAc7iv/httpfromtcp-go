package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	buffer := make([]byte, 8)

	go func() {
		defer f.Close()
		defer close(ch)

		var lineBuffer strings.Builder

		for {
			n, err := f.Read(buffer)
			if n > 0 {
				for _, b := range buffer[:n] {
					if b == '\n' {
						ch <- lineBuffer.String()
						lineBuffer.Reset()
					} else {
						lineBuffer.WriteByte(b)
					}
				}
			}

			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading file:", err)
				}
				if lineBuffer.Len() > 0 {
					ch <- lineBuffer.String()
				}
				break
			}
		}
	}()

	return ch
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")

	fmt.Println("TCP Server listening on port:42069")

	if err != nil {
		fmt.Printf("Error while listening: %v", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Print(err)
			break
		}

		defer conn.Close()
		fmt.Println("New Connection accepted!")

		readChannel := getLinesChannel(conn)

		for msg := range readChannel {
			fmt.Printf("%s\n", msg)
		}

		fmt.Println("Connection closed!")
	}

	defer listener.Close()
}
