package main

import (
	"fmt"
	"github.com/civet148/gotools/process"
	"os"
	"time"
)

func main() {
	strName := process.GetProcessName()
	strPath := process.GetProcessPath()
	pid := process.GetPID()
	strEnvPath := process.GetEnvPathSlice()

	if len(os.Args) > 1 && os.Args[1] == "--daemon" { //go run main.go --daemon
		switch os.Args[1] {
		case "--daemon":
			{
				fmt.Printf("running as daemon mode\n")
				process.Daemon()
				for i := 0; i < 1; i++ {
					time.Sleep(30 * time.Second) //sleep 30s and exit...
				}
			}
		}
	}
	fmt.Printf("process path [%v]\n", strPath)
	fmt.Printf("process name [%v]\n", strName)
	fmt.Printf("process pid [%v]\n", pid)
	fmt.Printf("process env %+v\n", strEnvPath)
}
