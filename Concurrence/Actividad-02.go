package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var (
		opc = ""
		condition = true
		flag = make(chan bool)
		scan = bufio.NewScanner(os.Stdin)

		processId uint64
		newProcess *Process
	 	processIdCreate uint64
	 	processIdDelete string

		process []*Process
		processAdmin = &ProcessManager {
			Processes: process,
		}
	)

	for exit := true; exit; exit = condition {
		fmt.Print("\033[H\033[2J")
		fmt.Println("******* Administrador de Procesos *******")
		fmt.Println("1) Agregar Proceso")
		fmt.Println("2) Mostrar Procesos")
		fmt.Println("3) Eliminar Proceso")
		fmt.Println("0) Salir")
		fmt.Print(" >> ")
		scan.Scan()
		opc = scan.Text()

		switch opc {
			case "1":
				newProcess = NewProcess(processIdCreate)
				processAdmin.AddProcess(newProcess)
				processIdCreate += 1
				go newProcess.Start()
				fmt.Print("\n-Proceso Añadido #", processIdCreate, "\n\n")
				time.Sleep(time.Millisecond * 2000)
				break
			case "2":
				if len(processAdmin.Processes) != 0 {
					go Concurrently(processAdmin, flag)
					scan.Scan()
					flag <- true
				} else {
					fmt.Print("\n-No hay procesos\n\n")
					time.Sleep(time.Millisecond * 3000)
				}
				break
			case "3":
				fmt.Print("\n-Dame el ID del proceso a eliminar: ")
				scan.Scan()
				processIdDelete = scan.Text()
				processId, _ = strconv.ParseUint(processIdDelete, 10, 64)
				if processAdmin.KillProcess(processId) {
					fmt.Print("-El Proceso no. ", processId, " fue eliminado.\n\n")
					time.Sleep(time.Millisecond * 2000)
				} else {
					fmt.Print("-No se encontro el proceso\n\n")
					time.Sleep(time.Millisecond * 3000)
				}
				processId = 0
				break
			case "0":
				processAdmin.exited()
				condition = false
				break
			default:
				invalidOptions()
		}
	}
}

type Process struct {
	Id uint64
	Task uint64
	IsRunning bool
}

func (process *Process) Start() {
	process.Task = 0
	process.IsRunning = true

	for {
		process.Task += 1
		time.Sleep(time.Millisecond * 500)
		if !process.IsRunning {
			break
		}
	}
}

func (process *Process) Stop() {
	process.IsRunning = false
}

func NewProcess(id uint64) *Process {
	return &Process{
		Id: id,
	}
}

type ProcessManager struct {
	Processes []*Process
}

func (processAdmin *ProcessManager) AddProcess(process *Process) {
	processAdmin.Processes = append(processAdmin.Processes, process)
}

func (processAdmin *ProcessManager) KillProcess(processId uint64) bool {
	var newProcess []*Process
	deleted := false

	for _, process := range processAdmin.Processes {
		if process.Id != processId {
			newProcess = append(newProcess, process)
		}

		if process.Id == processId {
			deleted = true
			process.Stop()
		}
	}

	processAdmin.Processes = newProcess
	return deleted
}

func (processAdmin *ProcessManager) ShowProcess() {
	for _, process := range processAdmin.Processes {
		fmt.Print("\nProceso #", process.Id, " : ", process.Task)
	}
	fmt.Println()
}

func Concurrently(processAdmin *ProcessManager, flag chan bool) {
	for {
		select {
		case <-flag:
			return
		default:
			processAdmin.ShowProcess()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func invalidOptions() {
	fmt.Print("\nOpción Invalida\n\n")
}

func (processAdmin *ProcessManager) exited() {
	for _, process := range processAdmin.Processes {
		process.Stop()
	}
	fmt.Println("\nGracias por usar el programa :)")
}