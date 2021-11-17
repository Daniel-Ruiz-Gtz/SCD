package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Student struct {
	Id       uint64             `json:"id"`
	Name     string             `json:"nombre"`
	Subjects map[string]float64 `json:"materias"`
}

var studentsGroup map[uint64]Student
var nextId int

func students(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path)
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		studentsJson, err := json.MarshalIndent(studentsGroup, "", "    ")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(res, string(studentsJson))
	case "POST":
		var student Student
		err := json.NewDecoder(req.Body).Decode(&student)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		studentsGroup[student.Id] = student
		fmt.Fprintf(res, "{\n    \"Mensaje\": \"estudiante registrado\"\n}")
	}
}

func studentsById(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path)
	res.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseUint(strings.TrimPrefix(req.URL.Path, "/estudiantes/"), 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	student, isValid := studentsGroup[id]
	if !isValid {
		http.Error(res, "{\n    \"Error\": \"estudiante no encontrado\"\n}", http.StatusNotFound)
		return
	}

	switch req.Method {
	case "GET":
		studentJson, err := json.MarshalIndent(student, "", "    ")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = req.ParseForm(); err != nil {
			fmt.Fprintf(res, "Error en ParseForm() %v", err)
			return
		}

		fmt.Fprintf(res, string(studentJson))
	case "DELETE":
		delete(studentsGroup, id)
		fmt.Fprintf(res, "{\n    \"Mensaje\": \"estudiante eliminado\"\n}")
	case "PUT":
		subject := req.FormValue("materia")
		grade, err := strconv.ParseFloat(req.FormValue("cal"), 64)
		if err != nil {
			http.Error(res, "{\n    \"Error\": \"calificación invalida\"\n}", http.StatusBadRequest)
			return
		}

		student.Subjects[subject] = grade
		fmt.Fprintf(res, "{\n    \"Mensaje\": \"calificación actualizada\"\n}")
	}
}

func main() {
	studentsGroup = make(map[uint64]Student)
	nextId = 1

	http.HandleFunc("/estudiantes", students)
	http.HandleFunc("/estudiantes/", studentsById)

	fmt.Println("Servidor Corriendo... :)")
	fmt.Println("http://localhost:9000/")
	http.ListenAndServe(":9000", nil)
}
