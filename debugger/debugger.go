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
		// TODO: Add conditions to break
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
	fmt.Println("yabfig debugger for Brainfuck.\nCommands are similar to gdb, type \"help\" for a list of compatible commands.")
	scanner := bufio.NewScanner(os.Stdin)
  dbg.init()
	for {
		fmt.Printf("(yabfig-dbg) ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)
		inputs := strings.Split(input, " ")
		if inputs[0] == "q" || inputs[0] == "quit" {
			var confirm string
			fmt.Printf("Are you sure you want to quit (y/n)? ")
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) == "y" {
				os.Exit(0)
			}
		} else if inputs[0] == "h" || inputs[0] == "help" {
			fmt.Println("List of commands:")
			fmt.Println()
			fmt.Println("run -- Run the program")
			fmt.Println("print [pos] -- Print value at memory position")
			fmt.Println("next [count] -- Execute next instruction[s]")
			fmt.Println("jump [pos] -- Jump to a program position and resume")
			fmt.Println("break [pos] -- Add breakpoint at program position")
			fmt.Println("clear [pos] -- Delete breakpoint at program position")
			fmt.Println("watch [n = x] -- Set watchpoint when memory position n is x")
			fmt.Println("kill -- Kill program execution")
		} else if inputs[0] == "r" || inputs[0] == "run" {
			if dbg.running {
				var confirm string
				fmt.Printf("Program is already running, do you want to start from the " +
					"beginning (y/n)? ")
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					continue
				}
			}
			fmt.Printf("Running program: %s\n", dbg.programPath)
			dbg.interpreter = &bf.Interpreter{}
			dbg.interpreter.LoadFromFile(dbg.programPath)
			dbg.running = true
			dbg.runClocks()
		} else if inputs[0] == "p" || inputs[0] == "print" {
			if dbg.running {
				if len(inputs) > 1 {
					inputPos, err := strconv.Atoi(inputs[1])
					if err != nil {
						fmt.Println("Position must be an integer!")
						continue
					}
					value := dbg.interpreter.GetProperMemoryValue(inputPos)
					fmt.Printf("$%d = %d (%x)\n", inputPos, value, value)
				} else {
					fmt.Println("Not enough arguments for this command!")
				}
			} else {
				fmt.Println("Program is not running!")
			}
		} else if inputs[0] == "j" || inputs[0] == "jump" {
			if dbg.running {
				if len(inputs) > 1 {
					inputPos, err := strconv.Atoi(inputs[1])
					if err != nil {
						fmt.Println("Position must be an integer!")
						continue
					}
					dbg.interpreter.ProgramPosition = inputPos
					dbg.runClocks()
				} else {
					fmt.Println("Not enough arguments for this command!")
				}
			} else {
				fmt.Println("Program is not running!")
			}
		} else if inputs[0] == "n" || inputs[0] == "next" {
			if dbg.running {
				if len(inputs) > 1 {
					count, err := strconv.Atoi(inputs[1])
					if err != nil {
						fmt.Println("Count must be an integer!")
						continue
					}
					for count > 0 && dbg.interpreter.Clock() {
						count--
					}
				} else {
					dbg.interpreter.Clock()
				}
				dbg.checkIfRunning()
			} else {
				fmt.Println("Program is not running!")
			}
		} else if inputs[0] == "c" || inputs[0] == "continue" {
			if dbg.running {
				dbg.runClocks()
			} else {
				fmt.Println("Program is not running!")
			}
		} else if inputs[0] == "b" || inputs[0] == "break" {
			if len(inputs) > 1 {
				point, err := strconv.Atoi(inputs[1])
				if err != nil {
					fmt.Println("Breakpoint must be an integer!")
					continue
				}
				b, ok := dbg.breakpoints[point]
				if ok && b {
					fmt.Printf("Breakpoint already exists at position %d\n", point)
				} else {
					dbg.breakpoints[point] = true
					fmt.Printf("Breakpoint #%d at position %d\n", len(dbg.breakpoints), point)
				}
			} else {
				fmt.Println("Not enough arguments for this command!")
			}
		} else if inputs[0] == "clear" {
			if len(inputs) > 1 {
				point, err := strconv.Atoi(inputs[1])
				if err != nil {
					fmt.Println("Breakpoint must be an integer!")
					continue
				}
				b, ok := dbg.breakpoints[point]
				if ok && b {
					dbg.breakpoints[point] = false
					fmt.Printf("Breakpoint cleared at position %d \n", point)
				} else {
					fmt.Printf("A breakpoint does not exist at position %d\n", point)
				}
			} else {
				fmt.Println("Not enough arguments for this command!")
			}
		} else if inputs[0] == "watch" {
			if len(inputs) > 1 {
				s := inputs[1:]
				if len(s) != 3 {
					fmt.Println("Invalid expression! 1")
					continue
				}

				n, err := strconv.Atoi(s[0])
				x, err2 := strconv.Atoi(s[2])
				if err != nil || err2 != nil {
					fmt.Println("Invalid expression! 2")
					continue
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
					continue
				}

				dbg.watchpoints[n] = expression{eq, x}

			} else {
				fmt.Println("Not enough arguments for this command!")
			}

		} else if inputs[0] == "kill" {
      if dbg.running {
        dbg.running = false
      } else {
        fmt.Println("Cannot kill a program which is not running!")
      }
		} else if inputs[0] == "" || len(inputs) == 0 {
			// No input, ignore
		} else {
			fmt.Printf("Undefined command: \"%s\".\tTry \"help\".\n", inputs[0])
		}
	}
}
