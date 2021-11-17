package main

// Import libraries.
import (
	"./data"
	"fmt"
	"net"
	"net/rpc"
)

// Const type and port server.
const (
	Type = "tcp"
	Port = "localhost:9999"
)

// Function connection and start the server.
func server() {
	/* The new built-in function allocates memory.
	 * The first argument is a type, not a value, and the value returned
	 * is a pointer to a newly allocated zero value of that type.
	 */
	ser := new(data.Server)

	// Added structs subjects and students a server.
	ser.Subjects = make(map[string]map[string]float64)
	ser.Students = make(map[string]map[string]float64)

	//Register publishes the receiver's methods in the DefaultServer.
	_ = rpc.Register(ser)

	// Listening server.
	ln, err := net.Listen(Type, Port)

	// Error.
	if err != nil {
		fmt.Println("ERROR_ESCUCHAR: ", err)
	}

	for {
		// Accept connection for client.
		con, err := ln.Accept()

		// Error.
		if err != nil {
			fmt.Println("ERROR_CONECTAR: ", err)
			continue
		}

		// Listens for the client's request.
		go rpc.ServeConn(con)
	}
}

// Function main.
func main() {
	// Start server.
	go server()

	// Stop server.
	var exitServer string
	_,_ = fmt.Scanln(&exitServer)
}
