package main

import (
	"fmt"
	"github.com/civet148/gotools/gopprof"
)

func main() {

	//open your browser and type http://127.0.0.1:8999/debug/pprof
	//thread create -> http://127.0.0.1:8999/debug/pprof/threadcreate?debug=1
	//goroutine -> http://127.0.0.1:8999/debug/pprof/goroutine?debug=1
	//heap -> http://127.0.0.1:8999/debug/pprof/heap?debug=1
	if err := gopprof.Start("127.0.0.1:8999", true); err != nil {
		fmt.Printf("%s", err)
		return
	}
}
