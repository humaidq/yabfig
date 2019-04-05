package main

import (
	"fmt"
	dbg "gitlab.com/humaid/yabfig/debugger"
	bf "gitlab.com/humaid/yabfig/interpreter"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-debug" {
			dbg := dbg.Debugger{}
			if len(args) > 1 {
				dbg.SetProgram(args[len(args)-1])
			}
			dbg.RunDebugger()
		} else {
			ipr := bf.Interpreter{}
			ipr.LoadFromFile(args[len(args)-1])
			if len(args) > 1 && args[0] == "-lint" {
				fmt.Printf("%s\n", ipr.Program)
			} else {
				ipr.Run()
			}
		}
	} else {
		fmt.Println("Usage: yabfig [option] <file>")
		fmt.Println("Options:")
		fmt.Println("\t-lint\t\tLint (format) a Brainfuck file by removing spaces and non-instruction characters and output it to standard output.")
		fmt.Println("\t-debug\t\tRun an interactive gdb-style debugger.")
	}
}
