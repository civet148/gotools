package main

import (
	"fmt"
	"github.com/civet148/gotools/process"
)

func main() {
	strName := process.GetProcessName()
	strPath := process.GetProcessPath()
	pid := process.GetPID()
	strEnvPath := process.GetEnvPathSlice()
	fmt.Printf("process path [%v] name [%v] pid [%v] env path %+v\n", strPath, strName, pid, strEnvPath)
}
