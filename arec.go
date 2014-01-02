package main

import (
	"fmt"
	"github.com/yosisa/arec/epg"
	"os"
)

func main() {
	data, err := epg.DecodeJson(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data[0].Name)
	for _, program := range data[0].Programs {
		fmt.Printf("%+v\n", program)
	}
}
