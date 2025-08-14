package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	defer f.Close()

	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("ERR: ", err)
		}
		fmt.Printf("read: %s\n", string(data[:n]))
	}
}
