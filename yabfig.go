package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		data, err := ioutil.ReadFile(args[len(args)-1])
		if err != nil {
			log.Fatal(err)
		}
		ipr := Interpreter{}
		ipr.Init()
		ipr.LoadProgram(data)
		if len(args) > 1 && args[0] == "-lint" {
			fmt.Printf("%s\n", ipr.program)
		} else {
			for ipr.Clock() {
			}
		}
	} else {
		fmt.Println("yabfig: usage: yabfig [-lint] [file]")
	}
}
