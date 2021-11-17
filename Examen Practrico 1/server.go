package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"./requests"
)

var printMessages = true

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("No se pudo iniciar el servidor: ", err.Error())
		return
	}
	defer ln.Close()

	clients := make(map[string]net.Conn)
	reqs := make([]requests.Request, 0)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("No se pudo completar la conexión: ", err.Error())
				continue
			}

			go handleClientRequests(conn, clients, &reqs)
		}
	}()

	input := -1
	fmt.Print("\033[H\033[2J")
	fmt.Println("**SERVIDOR**")
	fmt.Println("1) Respaldar mensajes")
	fmt.Println("0) Salir")
	for input != 0 {
		fmt.Scan(&input)

		switch input {
		case 1:
			fmt.Println("Respaldando mensajes...")
			createBackup(&reqs)
		case 0:
			fmt.Println("Apagando servidor...")
		default:
			fmt.Println("Opción invalida")
		}
	}
}

func handleClientRequests(conn net.Conn, clients map[string]net.Conn, reqs *[]requests.Request) {
	for {
		var req requests.Request
		err := gob.NewDecoder(conn).Decode(&req)
		if err != nil {
			fmt.Println("No se pudo obtener el request: ", err.Error())
			continue
		}

		switch req.Id {
		case requests.HELLO:
			clients[req.Sender] = conn
			if printMessages {
				fmt.Printf("%s se ha conectado.\n", req.Sender)
			}
		case requests.TXT:
			if printMessages {
				fmt.Printf("%s escribió: %s\n", req.Sender, req.Info)
			}
		case requests.FILE:
			if printMessages {
				fmt.Printf("%s envió el archivo: %s\n", req.Sender, req.Info)
			}
		case requests.GOODBYE:
			delete(clients, req.Sender)
			if printMessages {
				fmt.Printf("%s se ha desconectado.\n", req.Sender)
			}
			broadcastRequest(conn, clients, req)
			*reqs = append(*reqs, req)
			return
		}

		*reqs = append(*reqs, req)
		broadcastRequest(conn, clients, req)
	}

}

func broadcastRequest(conn net.Conn, clients map[string]net.Conn, request requests.Request) {
	for _, v := range clients {
		if conn != v {
			err := gob.NewEncoder(v).Encode(request)
			if err != nil {
				fmt.Println("No se pudo enviar el request ", err.Error())
				continue
			}
		}
	}
}

func createBackup(reqs *[]requests.Request) {
	if _, err := os.Stat("backup"); os.IsNotExist(err) {
		os.Mkdir("backup", 0755)
	}

	file, err := os.Create("./backup/respaldo.txt")
	if err != nil {
		fmt.Println("No se pudo crear el respaldo: ", err.Error())
		return
	}
	defer file.Close()

	for _, req := range *reqs {
		switch req.Id {
		case requests.HELLO:
			file.WriteString(req.Sender + " se ha conectado." + "\n")
		case requests.TXT:
			file.WriteString(req.Sender + " escribió: " + req.Info + "\n")
		case requests.FILE:
			file.WriteString(req.Sender + " envió el archivo: " + req.Info + "\n")
		case requests.GOODBYE:
			file.WriteString(req.Sender + " se ha desconectado." + "\n")
		}
	}
	fmt.Println("Respaldo completado.")
}
