package main

import (
	"artemisLogParser/logparser"
	"fmt"
	"os"
)

func main() {
	r, _ := os.Open("C:\\Users\\zivno\\Documents\\3d game design\\gameDesign\\GameJamJan21\\GameLogs\\log-2023-02-24T10.13.26.artemis")
	defer r.Close()
	fmt.Println(logparser.Read(r))
}
