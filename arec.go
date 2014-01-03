package main

import (
	"fmt"
	"github.com/yosisa/arec/epg"
	"github.com/yosisa/arec/reserve"
	"os"
)

func main() {
	reserve.Connect("mongodb://localhost/arec")

	data, err := epg.DecodeJson(os.Stdin)
	if err != nil {
		panic(err)
	}
	channel := data[0]
	if err := channel.Save(); err != nil {
		fmt.Printf("%+v\n", err)
	}

	for _, program := range channel.Programs {
		if err := program.Save(); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
