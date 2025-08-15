package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Failed to resolve UDP Addr", err)
	}

	// no need for a specific local addr; nil lets the OS choose
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalf("Failed to Dial UDP Addr: %s, Reason: %+v", raddr, err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		// Read a line (includes the trailing '\n' if present)
		line, err := reader.ReadString('\n')
		if err != nil {
			// Spec: log any read errors (don’t crash the program)
			log.Printf("ReadString error: %v", err)
			// Continue to allow more input attempts
			continue
		}

		// Write the line to the UDP connection
		if _, err := conn.Write([]byte(line)); err != nil {
			// Spec: log any write errors
			log.Printf("UDP write error: %v", err)
			continue
		}
	}
}
