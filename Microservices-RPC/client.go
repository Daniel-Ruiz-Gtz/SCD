package main

// Import libraries.
import (
	"fmt"
	"net/rpc"
	"time"
)

// Register -> data struct.
type Register struct {
	Student string
	Subject string
	Grade   float64
}

// Const type and port server.
const (
	ConnType = "tcp"
	ConnPort = "localhost:9999"
)

// Function connection in the server.
func client() {
	// Connection in the server.
	con, err := rpc.Dial(ConnType, ConnPort)

	// Error.
	if err != nil {
		fmt.Println("ERROR_CONNECTING: ", err)
		return
	}

	// Variable option for menu.
	var option int64
	var resultDone string
	var resultAverage float64

	// Instance register.
	var data Register

	for {
		// Menu.
		fmt.Print("\033[H\033[2J")
		fmt.Println("\n******* Microservicio  (RPC) *******\n")
		fmt.Println("1) Agregar la calificación de un alumno por materia")
		fmt.Println("2) Obtener el promedio del alumno")
		fmt.Println("3) Obtener el promedio de todos los alumnos")
		fmt.Println("4) Obtener el promedio por materia")
		fmt.Println("0) Salir del Sistema")
		// Get option.
		fmt.Print(" >> ")
		_, _ = fmt.Scanln(&option)

		switch option {
			case 1:
				registerData(data, con, resultDone) // Register data.
			case 2:
				studentAvg(data, con, resultAverage) // Student average.
			case 3:
				generalAvg(data, con, resultAverage) // General average.
			case 4:
				subjectAvg(data, con, resultAverage) // Subject average.
			case 0:
				exited() // Exit and terminate process.
				return
			default:
				invalidOptions() // Invalid Options.
			}
	}
}

/* Function registerData
	@param data Register
	@param con *rpc.Client
	@param resultDone string
*/
func registerData(data Register, con *rpc.Client, resultDone string) {
	// Register data.
	fmt.Println("\n******* NUEVO ITEM *******")
	fmt.Print("- Nombre: ")
	_, _ = fmt.Scanln(&data.Student)
	fmt.Print("- Materia: ")
	_, _ = fmt.Scanln(&data.Subject)
	fmt.Print("- Calificación: ")
	_, _ = fmt.Scanln(&data.Grade)

	// To call a method exposed by RPC we use the Call() method.
	err := con.Call("Server.AddRegister", data, &resultDone)

	// Error.
	if err != nil {
		fmt.Println("ERROR_ENVIAR_INFO", err)
	}
}

/* Function studentAvg
	@param data Register
	@param con *rpc.Client
	@param resultAverage float64
*/
func studentAvg(data Register, con *rpc.Client, resultAverage float64) {
	// Get name for search.
	fmt.Println("\n******* PROMEDIO *******")
	fmt.Print("Nombre: ")
	_, _ = fmt.Scanln(&data.Student)

	// To call a method exposed by RPC we use the Call() method.
	err := con.Call("Server.StudentAverage", data, &resultAverage)

	// Error.
	if err != nil {
		fmt.Println("ERROR_OBTENER_INFO", err)
	} else {
		fmt.Println("Promedio ", data.Student, ": ", resultAverage)
		time.Sleep(4 * time.Second)
	}
}

/* Function generalAvg
	@param data Register
	@param con *rpc.Client
	@param resultAverage float64
*/
func generalAvg(data Register, con *rpc.Client, resultAverage float64) {
	// To call a method exposed by RPC we use the Call() method.
	fmt.Println("\n******* PROMEDIO GENERAL *******")
	err := con.Call("Server.GeneralAverage", data, &resultAverage)

	// Error.
	if err != nil {
		fmt.Println("ERROR_OBETNER_INFO", err)
	} else {
		fmt.Println("Promedio: ", resultAverage)
		time.Sleep(4 * time.Second)
	}
}

/* Function subjectAvg
	@param data Register
	@param con *rpc.Client
	@param resultAverage float64
*/
func subjectAvg(data Register, con *rpc.Client, resultAverage float64) {
	// Get subject for search.
	fmt.Println("\n******* PROMEDIO POR MATERIA *******")
	fmt.Print("Materia: ")
	_, _ = fmt.Scanln(&data.Subject)

	// To call a method exposed by RPC we use the Call() method.
	err := con.Call("Server.SubjectAverage", data, &resultAverage)

	// Error.
	if err != nil {
		fmt.Println("ERROR_OBTENER_INFO", err)
	} else {
		fmt.Println("EL promedio de la materia de ", data.Subject, " es de: ", resultAverage)
		time.Sleep(4 * time.Second)
	}
}

// Function invalidOptions.
func invalidOptions() {
	fmt.Print("\nOpción Invalida\n\n")
}

// Function exited.
func exited() {
	fmt.Println("\nGracias por usar el programa :)")
}

// Function main.
func main() {
	// Start client.
	client()
}