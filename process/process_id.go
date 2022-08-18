package process

import "os"

func GetPID() int {
	return os.Getpid()
}

func GetGID() int {
	return os.Getgid()
}

func GetUID() int {
	return os.Getuid()
}

func GetParentID() int {
	return os.Getppid()
}
