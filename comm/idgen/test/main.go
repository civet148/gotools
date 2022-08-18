package main

import (
	"fmt"
	"github.com/civet148/gotools/comm/idgen"
)

func main() {
	var nodeId, businessId int64
	nodeId = 1     //host node id
	businessId = 1 //business id
	workId := nodeId + businessId
	if g, err := idgen.NewIdGenerator(workId); err != nil {
		fmt.Println(err.Error())
	} else {
		for i := 0; i < 10; i++ {
			fmt.Println(fmt.Sprintf("[%d] id %d", i, g.CreateId()))
		}
	}
}
