package main

// Import libraries.
import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

// Const type and port server.
const (
	Type = "tcp"
	Port = "localhost:9999"
)

// Function connection and start the server.
func server(channels []chan uint64, exit []chan bool) {
	// added data, send data and receive data.
	add := make(chan uint64, 2)
	send := make(chan net.Conn)
	receive := make(chan bool)

	// Listening server.
	listener, err := net.Listen(Type, Port)

	if err != nil {
		log.Fatalln("ERROR_LISTENING: ", err)
	}

	// Closer connection.
	defer listener.Close()

	// If you want, you can increment a counter here and inject to ProcessManager below as client identifier
	go processManager(channels, exit, add, send, receive)

	// Accept connection for client.
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println("ERROR_CONNECTING: ", err)
			continue
		}
		send <- con
	}
}

// Function handleClientRequest -> arguments of connection, channels of add, recive and exit.
func handleClientRequest(con net.Conn, channel chan uint64, exit chan bool, add chan uint64, receive chan bool) {
	// Terminate connection.
	defer con.Close()

	// Terminate request of client.
	exit <- true

	// Data
	data := [2]uint64{<-channel, <- channel}
	err := gob.NewEncoder(con).Encode(data)

	if err != nil {
		log.Println("ERROR: ", err)
	}

	for {
		// Waiting for the client request
		err := gob.NewDecoder(con).Decode(&data[1])

		if err != nil {
			add <- data[0]
			add <- data[1]
			receive <- true
			return
		}
	}
}

// Function process -> increment for process and select action for run.
func process(id uint64, i uint64, channel chan uint64, exit chan bool) {
	for {
		select {
		case <-exit:
			channel <- id
			channel <- i
			return
		default:
			fmt.Println("ID ", id, " : ", i)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// Function processManager -> receive or send process.
func processManager(channels []chan uint64, exit []chan bool, add chan uint64, send chan net.Conn, receive chan bool) {
	for {
		select {
		case <- receive:
			inf := [2]uint64{<-add, <- add}
			channels = append(channels, make(chan uint64, 2))
			exit = append(exit, make(chan bool))
			go process(inf[0], inf[1], channels[len(channels)-1], exit[len(exit)-1])
		case c := <- send:
			go handleClientRequest(c, channels[0], exit[0], add, receive)
			channels = channels[1:]
			exit = exit[1:]
		}
	}
}

// Function main.
func main() {
	// variables
	var idCount uint64 = 0
	channels := make([]chan uint64, 5)
	exit := make([]chan bool, 5)

	// Quantity of process.
	for i:= 0; i < 5; i++{
		channels[i] = make(chan uint64, 2)
		exit[i] = make(chan bool)
		go process(idCount, 0, channels[i], exit[i])
		idCount++
	}

	// Start server.
	go server(channels, exit)

	// Stop server.
	var exitServer string
	_,_ = fmt.Scanln(&exitServer)
}
