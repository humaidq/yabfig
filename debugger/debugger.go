package debugger

import (
	"bufio"
	"fmt"
	bf "gitlab.com/humaid/yabfig/interpreter"
	"os"
	"strconv"
	"strings"
)

const (
	equal = iota
	notEqual
	greater
	less
	greaterEqual
	lessEqual
)

type equality int

var equalityNames = []string{
	"=",
	"!=",
	">",
	"<",
	">=",
	"<=",
}

type expression struct {
	equal equality
	value int
}

// Debugger represents a debugging session.
type Debugger struct {
	programPath string
	interpreter *bf.Interpreter
	running     bool
	breakpoints map[int]bool
	watchpoints map[int]expression
}

// SetProgram sets the program path for the debugger.
func (dbg *Debugger) SetProgram(program string) {
	dbg.programPath = program
}

func (dbg *Debugger) runClocks() {
	for dbg.interpreter.Clock() {
		pos := dbg.interpreter.ProgramPosition
		b, ok := dbg.breakpoints[pos]
		if ok && b {
			fmt.Printf("Breakpoint hit at position %d\n", pos)
			return
		}

		for k, v := range dbg.watchpoints {
			mem := dbg.interpreter.GetProperMemoryValue(k)
			match := false
			switch v.equal {
			case equal:
				match = (mem == v.value)
			case notEqual:
				match = (mem != v.value)
			case greater:
				match = (mem > v.value)
			case less:
				match = (mem < v.value)
			case greaterEqual:
				match = (mem >= v.value)
			case lessEqual:
				match = (mem <= v.value)
			}

			if match {
				fmt.Printf("Watchpoint hit at position %d (%d%s%d)\n", pos, mem,
					equalityNames[v.equal], v.value)
				return
			}
		}

	}
	dbg.checkIfRunning()
}

func (dbg *Debugger) checkIfRunning() {
	if dbg.interpreter != nil && dbg.running && dbg.interpreter.IsEnded() {
		fmt.Println("Program exited")
		dbg.running = false
	}
}

func (dbg *Debugger) init() {
	dbg.breakpoints = make(map[int]bool)
	dbg.watchpoints = make(map[int]expression)
}

type command struct {
	commandName        []string
	argumentLabels     string
	commandDescription string
	minArguments       int
	requireRunning     bool
	function           func(*Debugger, []string)
}

var commands = []command{
	command{[]string{"run", "r"}, "", "Run the program", 0, false, run},
	command{[]string{"quit", "q"}, "", "Quit the debugger", 0, false, quit},
	command{[]string{"file", "f"}, "[path]", "Load a program from a file path", 0, false, file},
	command{[]string{"print", "p"}, "[pos]", "Print value at memory position", 1, true, printVal},
	command{[]string{"continue", "cont", "c"}, "", "Continue execution", 0, true, cont},
	command{[]string{"next", "n"}, "[count]", "Execute next instruction[s]", 0, true, next},
	command{[]string{"jump", "j"}, "[pos]", "Jump to a program position and resume", 1, true, jump},
	command{[]string{"break", "b"}, "[pos]", "Add a breakpoint at program position", 1, false, addBreakpoint},
	command{[]string{"clear", "cl"}, "[pos]", "Clear breakpoint at program position", 1, false, clear},
	command{[]string{"watch", "w"}, "[n = x]", "Set watchpoint when memory position n is x", 3, false, addWatchpoint},
	command{[]string{"kill"}, "", "Stop the program", 0, true, kill},
}

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

// RunDebugger starts the interactive debugger.
func (dbg *Debugger) RunDebugger() {
	fmt.Println("yabfig debugger for Brainfuck.\nCommands are similar to gdb, type \"help\" for a list of compatible commands.")
	scanner := bufio.NewScanner(os.Stdin)
	dbg.init()
	for {
		fmt.Printf("(yabfig-dbg) ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)
		inputs := strings.Split(input, " ")

		cmd := getCommand(inputs[0])

		if inputs[0] == "" || len(inputs) == 0 {
			// No input, ignore
			continue
		}

		if inputs[0] == "help" || inputs[0] == "h" {
			fmt.Println("List of commands:")
			fmt.Println()
			for _, cmd := range commands {
				args := cmd.argumentLabels
				if len(args) > 0 {
					args = " " + args
				}
				fmt.Printf("%s%s \t\t %s\n", cmd.commandName[0], args, cmd.commandDescription)
			}

		}

		if cmd == nil {
			fmt.Printf("Undefined command: \"%s\".\tTry \"help\".\n", inputs[0])
			continue
		}

		if cmd.requireRunning && !dbg.running {
			fmt.Println("Program is not running!")
			continue
		}

		if len(inputs[1:]) < cmd.minArguments {
			fmt.Println("Not enough arguments for this command!")
			fmt.Printf("usage: %s %s\n", cmd.commandName[0], cmd.argumentLabels)
			continue
		}

		cmd.function(dbg, inputs[1:])
	}
}
