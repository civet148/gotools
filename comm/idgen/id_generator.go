package idgen

import (
	"fmt"
	gsf "github.com/zheng-ji/goSnowFlake"
)

type idGenerator struct {
	iw *gsf.IdWorker
}

func init() {
}

// workerId, 0 < workerId < 1024
func NewIdGenerator(workerId int64) (g *idGenerator, err error) {
	g = new(idGenerator)
	if workerId <= 0 || workerId > 1023 {
		err = fmt.Errorf("work id must > 0 and < 1023")
		fmt.Println("work id must > 0 and < 1023")
		return
	}
	// Params: Given the workerId, 0 < workerId < 1024
	g.iw, err = gsf.NewIdWorker(workerId)
	if err != nil {
		fmt.Println(fmt.Sprintf("snowflake new id worker error: %s", err.Error()))
		return
	}
	return
}

func (g *idGenerator) CreateId() (id int64) {
	var err error
	if id, err = g.iw.NextId(); err != nil {
		fmt.Println(err.Error())
	}
	return
}
