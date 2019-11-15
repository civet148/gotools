package idgen

import "testing"
import "fmt"

func Test_CreateId(t *testing.T) {
	id, err := CreateId()
	fmt.Println("err: ", err)
	fmt.Println("id: ", id)
}
