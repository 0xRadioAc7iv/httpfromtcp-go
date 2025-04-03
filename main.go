package main

import (
	"fmt"
	"io"
	"os"
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
	file, err := os.Open("messages.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	fileChannel := getLinesChannel(file)

	for msg := range fileChannel {
		fmt.Printf("read: %s\n", msg)
	}
}
