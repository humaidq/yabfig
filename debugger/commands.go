package debugger

import (
	"fmt"
	bf "gitlab.com/humaid/yabfig/interpreter"
	"os"
	"strconv"
	"strings"
)

func quit(dbg *Debugger, _ []string) {
	var confirm string
	fmt.Printf("Are you sure you want to quit (y/n)? ")
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) == "y" {
		os.Exit(0)
	}
}

func run(dbg *Debugger, _ []string) {
	if dbg.running {
		var confirm string
		fmt.Printf("Program is already running, do you want to start from the " +
			"beginning (y/n)? ")
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			return
		}
	}
	if len(dbg.programPath) == 0 {
		fmt.Println("No program loaded! Load with `file [path]`")
		return
	}
	fmt.Printf("Running program: %s\n", dbg.programPath)
	dbg.interpreter = &bf.Interpreter{}
	dbg.interpreter.LoadFromFile(dbg.programPath)
	dbg.running = true
	dbg.runClocks()
}

func printVal(dbg *Debugger, args []string) {
	inputPos, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Position must be an integer!")
		return
	}
	value := dbg.interpreter.GetProperMemoryValue(inputPos)
	fmt.Printf("$%d = %d (0x%08x)\n", inputPos, value, value)
}

func file(dbg *Debugger, args []string) {
	// TODO file check
	dbg.SetProgram(args[0])
}

func jump(dbg *Debugger, args []string) {
	inputPos, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Position must be an integer!")
		return
	}
	dbg.interpreter.ProgramPosition = inputPos
	dbg.runClocks()
}

func cont(dbg *Debugger, args []string) {
	dbg.runClocks()
}

func next(dbg *Debugger, args []string) {
	if len(args) > 0 {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Count must be an integer!")
			return
		}
		for count > 0 && dbg.interpreter.Clock() {
			count--
		}
	} else {
		dbg.interpreter.Clock()
	}
}

func addBreakpoint(dbg *Debugger, args []string) {
	point, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Breakpoint must be an integer!")
		return
	}
	b, ok := dbg.breakpoints[point]
	if ok && b {
		fmt.Printf("Breakpoint already exists at position %d\n", point)
	} else {
		dbg.breakpoints[point] = true
		fmt.Printf("Breakpoint #%d at position %d\n", len(dbg.breakpoints), point)
	}
}

func clear(dbg *Debugger, args []string) {
	point, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Breakpoint must be an integer!")
		return
	}
	b, ok := dbg.breakpoints[point]
	if ok && b {
		dbg.breakpoints[point] = false
		fmt.Printf("Breakpoint cleared at position %d \n", point)
	} else {
		fmt.Printf("A breakpoint does not exist at position %d\n", point)
	}
}

func addWatchpoint(dbg *Debugger, args []string) {
	s := args
	if len(s) != 3 {
		fmt.Println("Invalid expression! 1")
		return
	}

	n, err := strconv.Atoi(s[0])
	x, err2 := strconv.Atoi(s[2])
	if err != nil || err2 != nil {
		fmt.Println("Invalid expression! 2")
		return
	}
	var eq equality
	switch s[1] {
	case "=":
		eq = equal
	case "!=":
		eq = notEqual
	case ">":
		eq = greater
	case "<":
		eq = less
	case ">=":
		eq = greaterEqual
	case "<=":
		eq = lessEqual
	default:
		fmt.Println("Invalid equality sign!")
		return
	}

	dbg.watchpoints[n] = expression{eq, x}
}

func kill(dbg *Debugger, args []string) {
	dbg.running = false
}

func getCommand(input string) *command {
	for _, cmd := range commands {
		for _, alias := range cmd.commandName {
			if input == alias {
				return &cmd
			}
		}

	}
	return nil
}
