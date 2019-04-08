package debugger

import (
	"bufio"
	"fmt"
	bf "gitlab.com/humaid/yabfig/interpreter"
	"os"
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

// RunDebugger starts the interactive debugger.
func (dbg *Debugger) RunDebugger() {
	fmt.Println("yabfig debugger for BF.\nCommands are similar to gdb, type \"help\" for a list of compatible commands.")
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
      continue
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
