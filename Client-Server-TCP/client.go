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
	ConnType = "tcp"
	ConnPort = "localhost:9999"
)

// function connection in the server.
func client() {
	// Connection in the server.
	con, err := net.Dial(ConnType, ConnPort)

	if err != nil {
		log.Println("ERROR_CONNECTING: ", err)
	}

	// Terminated connection the client.
	defer con.Close()

	// Data.
	var dataClient [2]uint64
	err = gob.NewDecoder(con).Decode(&dataClient)

	if err != nil {
		log.Println(err)
	} else {
		// Run process.
		channel := make(chan uint64)
		go processClient(dataClient[0], dataClient[1], channel)

		for {
			dataClient[1] = <- channel
			err := gob.NewEncoder(con).Encode(dataClient[1])

			if err != nil {
				log.Println("ERROR: ", err)
				return
			}
		}
	}
}

// Function processClient -> run process in the client.
func processClient(id uint64, i uint64, channel chan uint64) {
	for {
		fmt.Println(id, " : ", i)
		i = i + 1
		channel <- i
		time.Sleep(time.Millisecond * 500)
	}
}

// function main.
func main() {
	// Start client.
	go client()

	// Stop client.
	var exitClient string
	_,_ = fmt.Scanln(&exitClient)
}
