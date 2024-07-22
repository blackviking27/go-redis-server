package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")

	// create a server
	l, err := net.Listen("tcp", ":6379")
	if (err) != nil {
		fmt.Println("Unable to start a server:")
		fmt.Println(err)
		return
	}

	// create a connection
	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error while accepting connections")
		fmt.Println(err)
		return
	}

	defer conn.Close()

	// create a infinte for loop to accept the commands on tcp

	for {
		buffer := make([]byte, 1024)

		// read data from client
		_, err := conn.Read(buffer)
		if err != nil {

			// if end of line is encouontered
			if err == io.EOF {
				fmt.Println("EOF reached")
				break
			}

			fmt.Println("Error while reading data from client", err.Error())
			os.Exit(1)
		}

		fmt.Println("Sending response back", string(buffer[:]))
		// send a response back to client
		conn.Write([]byte("+OK\r\n"))

	}
	fmt.Println("Shutting down the server")

}
