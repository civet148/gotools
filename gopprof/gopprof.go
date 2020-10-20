package gopprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

//strAddress -> http listen address
//note: if block=true this function will block your goroutine
//open your browser and type http://host:port/debug/pprof
//thread create -> http://host:port/debug/pprof/threadcreate?debug=1
//goroutine -> http://host:port/debug/pprof/goroutine?debug=1
//heap -> http://host:port/debug/pprof/heap?debug=1
func Start(strHttpAddr string, block bool) (err error) {

	if !block {
		go listen(strHttpAddr)
	} else {
		err = listen(strHttpAddr)
	}
	return
}

func listen(strHttpAddr string) (err error) {
	if err = http.ListenAndServe(strHttpAddr, nil); err != nil {
		fmt.Printf("listen %s error [%s]\n", strHttpAddr, err)
		return
	}
	return
}
