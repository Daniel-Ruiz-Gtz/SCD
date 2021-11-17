package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Server struct {
	messages []Message
}

func (server *Server) Handler(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s recibio un request de tipo: %s\n", req.URL.Path, req.Method)
	switch req.Method {
	case "GET":
		response, err := server.GetMessages()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(response)
	case "POST":
		var msg Message
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response := server.AddMessage(msg)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(response)
	}
}

func (server *Server) GetMessages() ([]byte, error) {
	jsonData, err := json.MarshalIndent(server.messages, "", "    ")
	if err != nil {
		return jsonData, nil
	}
	return jsonData, err
}

func (server *Server) AddMessage(message Message) []byte {
	fmt.Printf("Remitente: %s Contenido: %s\n", message.Sender, message.Content)
	server.messages = append(server.messages, message)
	fmt.Println(len(server.messages))
	return []byte(`{"code": "OK"}`)
}

func main() {
	var topic string
	fmt.Print("Ingresa el tema de la sala: ")
	fmt.Scanln(&topic)
	var addr string
	fmt.Print("Ingresa la direccion de la sala: ")
	fmt.Scanln(&addr)
	MakeRoom(topic, addr)
}

func NewServer() *Server {
	var s Server
	return &s
}

func MakeRoom(topic string, addr string) {
	http.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			LoadHTML("../static/room/index.html"),
			"Sala de "+topic,
			"Tema de la sala: "+topic,
		)
	})
	http.HandleFunc("/chat", NewServer().Handler)
	fmt.Printf("Sala de chat con tema: %s esta en linea en: %s\n", topic, addr)
	print(addr)
	http.ListenAndServe(addr, nil)

}

func LoadHTML(filename string) string {
	html, _ := ioutil.ReadFile(filename)
	return string(html)
}
