package main

import (
	"fmt"
	"net"
	"strings"
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
		// buffer := make([]byte, 1024)

		// read data from client
		// _, err := conn.Read(buffer)
		// if err != nil {

		// 	// if end of line is encouontered
		// 	if err == io.EOF {
		// 		fmt.Println("EOF reached")
		// 		break
		// 	}

		// 	fmt.Println("Error while reading data from client", err.Error())
		// 	os.Exit(1)
		// }

		// fmt.Println("Sending response back", string(buffer[:]))
		// // send a response back to client
		// conn.Write([]byte("+OK\r\n"))

		// Read commands from user using the NewResp object
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Encountered EOF while reading the input")
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// command received from the user
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		fmt.Println("Response received")
		fmt.Println(value)

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		result := handler(args)
		writer.Write(result)
	}

}
