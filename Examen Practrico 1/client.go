package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"./requests"
)

var printMessages = true

func main() {
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println("No se pudo establecer conexión al servidor: ", err.Error())
		return
	}
	defer conn.Close()

	name := readName()
	sendHello(conn, name)

	go handleServerRequests(conn, name)

	input := -1
	fmt.Print("\033[H\033[2J")
	fmt.Println("**Bienvenido ",name,"**")
	fmt.Println("1) Mandar mensaje")
	fmt.Println("2) Mandar archivo")
	fmt.Println("0) Salir")
	for input != 4 {

		fmt.Scan(&input)

		switch input {
		case 1:
			sendText(conn, name)
		case 2:
			sendFile(conn, name)
		case 0:
			fmt.Println("Hasta luego ", name,"!")
			sendGoobye(conn, name)
		default:
			fmt.Println("Opción invalida")
		}
	}
}

func readName() string {
	fmt.Print("Ingresa el nombre del cliente: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()
	return name
}

func handleServerRequests(conn net.Conn, name string) {
	for {
		var req requests.Request
		err := gob.NewDecoder(conn).Decode(&req)
		if err != nil {
			fmt.Println("No se pudo obtener el request: ", err.Error())
			continue
		}

		switch req.Id {
		case requests.HELLO:
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
			writeFile(name, req.Info, req.Data)
		case requests.GOODBYE:
			if printMessages {
				fmt.Printf("%s se ha desconectado.\n", req.Sender)
			}
		}
	}
}

func sendHello(conn net.Conn, sender string) {
	err := gob.NewEncoder(conn).Encode(requests.Request{Id: requests.HELLO, Sender: sender})
	if err != nil {
		fmt.Println("No se pudo enviar el HELLO Request: ", err.Error())
		return
	}
}

func sendText(conn net.Conn, sender string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	msg := scanner.Text()
	err := gob.NewEncoder(conn).Encode(requests.Request{Id: requests.TXT, Sender: sender, Info: msg})
	if err != nil {
		fmt.Println("No se pudo enviar el TXT Request: ", err.Error())
		return
	}
}

func sendFile(conn net.Conn, sender string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filename := scanner.Text()
	file := readFile(filename)
	if file == nil {
		fmt.Println("No se pudo leer el archivo: ", file)
		return
	}
	err := gob.NewEncoder(conn).Encode(requests.Request{Id: requests.FILE, Sender: sender, Info: filepath.Base(filename), Data: file})
	if err != nil {
		fmt.Println("No se pudo enviar el FILE Request: ", err.Error())
		return
	}
}

func sendGoobye(conn net.Conn, sender string) {
	err := gob.NewEncoder(conn).Encode(requests.Request{Id: requests.GOODBYE, Sender: sender})
	if err != nil {
		fmt.Println("No se pudo enviar el GOODBYE Request: ", err.Error())
		return
	}
}

func readFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("No se pudo leer el archivo: ", err.Error())
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("No se pudo leer la estructura del archivo: ", err.Error())
		return nil
	}

	size := stat.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("No se pudo leer datos del archivo: ", err.Error())
		return nil
	}

	return buffer
}

func writeFile(receiver string, filename string, data []byte) {
	if _, err := os.Stat("received_files/" + receiver); os.IsNotExist(err) {
		os.MkdirAll("received_files/"+receiver, 0755)
	}

	file, err := os.Create("received_files/" + receiver + "/" + filename)
	if err != nil {
		fmt.Println("No se pudo crear el archivo: ", err.Error())
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("No se pudieron escribir datos en el archivo: ", err.Error())
	}
}
