package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func banProcess(port string, processNames []string, wg *sync.WaitGroup, stopChan chan string) {
	println("Checking port: ", port, " for process: ", processNames)
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			fmt.Printf("Goroutine for port %s stopped\n", port)
			return
		default:
			netCmd := exec.Command("netstat", "-ano")
			out, err := netCmd.Output()
			//println(string(out))
			if err != nil {
				println(err.Error())
			}

			lines := strings.Split(string(out), "\n")
			//println(lines[4])
			//println(processID)
			processId := ""
			for _, line := range lines[4:] {
				if strings.Contains(line, "0.0.0.0:"+port) {
					fields := strings.Fields(line)
					processId = fields[len(fields)-1]
					break
				}
			}
			if processId != "" {
				process := exec.Command("tasklist", "/fi", "PID eq "+processId)
				out, err := process.Output()
				//println(string(out))
				if err != nil {
					println(err.Error())
				}
				processName := ""
				for _, line := range strings.Split(string(out), "\n") {
					for _, process := range processNames {
						if strings.Contains(line, process) {
							processName = strings.Fields(line)[0]
						}
					}
				}
				//println("Process Name: ", processName)

				if strings.Contains(processName, "svchost") {
					_, err := exec.Command("taskkill", "/F", "/PID", processId).Output()
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Println("Port", port, "is in use. Killing the process..."+processName)
				}
			}
		}
	}
}

func banListProcess(ports []string, process []string) {
	var wg sync.WaitGroup
	stopChan := make(chan string)
	wg.Add(len(ports))
	for _, port := range ports {
		go banProcess(port, process, &wg, stopChan)
	}
	func() {
		println("To stop the process press enter N and press enter")
		var input string
		fmt.Scanln(&input)
		if input == "N" {
			stopChan <- "stop"
		}
	}()
	go func() {
		wg.Wait()
		close(stopChan)
	}()
	println("Done")
}
func main() {
	app := cli.App{
		Name:  "SvcHost Killer",
		Usage: "Kills the svchost process that is using the port 5555",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Value: "5555",
				Usage: "Port to check",
			},
			&cli.StringFlag{
				Name:  "ban",
				Value: "svchost",
				Usage: "Process to ban",
			},
		},
		Action: func(context *cli.Context) error {
			ports := context.String("port")
			process := context.String("ban")
			var portList, processList []string
			if ports != "" || process != "" {
				if strings.ContainsAny(ports, ",") {
					portList = strings.Split(ports, ",")
				} else {
					portList = append(portList, ports)
				}
				if strings.ContainsAny(process, ",") {
					processList = strings.Split(process, ",")
				} else {
					processList = append(processList, process)
				}
				println("Ports: ", portList)
				println("Process: ", processList)
				banListProcess(portList, processList)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
