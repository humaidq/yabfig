package main

import (
	"fmt"
	bf "gitlab.com/humaid/yabfig/interpreter"
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
		ipr := bf.Interpreter{}
		ipr.Init()
		ipr.LoadProgram(data)
		if len(args) > 1 && args[0] == "-lint" {
			fmt.Printf("%s\n", ipr.Program)
		} else {
			for ipr.Clock() {
			}
		}
	} else {
		fmt.Println("Usage: yabfig [options] <file>")
		fmt.Println("Options:")
		fmt.Println("\t-lint\t\tLint (format) a Brainf--k file by removing spaces and non-instruction characters and output it to standard output.")
	}
}
