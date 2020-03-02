package main

import (
	"github.com/civet148/gotools/comm/reflector"
	"github.com/civet148/gotools/log"
)

func main() {

	ParseStruct()
}

func ParseStruct() {
	type Person struct {
		Name    string  `json:"name" db:"name"`
		age     int     `json:"age" db:"age"` //not a export var
		Country string  `json:"country" db:"country"`
		Tall    float64 `json:"tall" db:"tall"`
	}

	var person = Person{
		Name:    "lory",
		age:     36,
		Country: "China",
		Tall:    1.75,
	}

	log.Debugf("parse to map %+v", reflector.Struct(&person).ToMap("json"))
}
