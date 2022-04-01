package grace

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func GracefullyExit(callback func(signal os.Signal) bool, signals ...os.Signal) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover [%+v]\n", r)
		}
	}()
	//capture signal of Ctrl+C and gracefully exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, signals...)
	go func() {
		for {
			select {
			case s := <-sig:
				if s != nil {
					fmt.Printf("capture signal [%v]\n", s.String())
					if callback == nil {
						if s == os.Interrupt {
							fmt.Printf("callback is nil and signal is %s, program exit...\n", s.String())
							exit0(sig)
						}
					} else {
						if !callback(s) {
							fmt.Printf("callback return false, program exit...\n")
							exit0(sig)
						}
					}
				}
			}
		}
	}()
}

func exit0(sig chan os.Signal) {
	time.Sleep(300 * time.Millisecond)
	close(sig)
	os.Exit(0)
}
