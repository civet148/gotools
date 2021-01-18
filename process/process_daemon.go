package process

import (
	"log"
	"os"
	"runtime"
	"syscall"
)

//return code, 0 success other means failure
func Daemon() (code int) {
	var r1, r2 uintptr
	var err syscall.Errno

	darwin := runtime.GOOS == "darwin"
	windows := runtime.GOOS == "windows"

	if windows {
		return 0
	}
	// already a daemon
	if syscall.Getppid() == 1 {
		return 0
	}
	// fork off the parent process
	r1, r2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return -1
	}
	// failure
	if r2 < 0 {
		os.Exit(-1)
	}
	// handle exception for darwin
	if darwin && r2 == 1 {
		r1 = 0
	}
	// if we got a good PID, then we call exit the parent process.
	if r1 > 0 {
		os.Exit(0)
	}
	/* Change the file mode mask */
	_ = syscall.Umask(0)
	// create a new SID for the child process
	sidRet, sidErrNo := syscall.Setsid()
	if sidErrNo != nil {
		log.Printf("Error: syscall.Setsid errno: %d\n", sidErrNo)
	}
	if sidRet < 0 {
		return -1
	}
	//redirect stdout/stdin/stderr to /dev/null
	f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if e == nil {
		fd := f.Fd()
		syscall.Dup2(int(fd), int(os.Stdin.Fd()))
		syscall.Dup2(int(fd), int(os.Stdout.Fd()))
		syscall.Dup2(int(fd), int(os.Stderr.Fd()))
	}
	return 0
}
