package main

import (
	"fmt"
	"github.com/zivoy/ArtemisLogParser/logparser"
	"os"
	"strings"
)

// todo write profiler cases

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("File path required")
		return
	}

	r, err := os.Open(strings.Join(args, " "))
	if err != nil {
		fmt.Println(err)
		return
	}
	g, err := logparser.Read(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(g)
	err = r.Close()
	if err != nil {
		println(err)
	}
}
