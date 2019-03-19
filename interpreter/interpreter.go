package interpreter

import (
	"fmt"
	"log"
	"strings"
)

const memorySize int = 500

type stack []int

func (s stack) push(v int) stack {
	return append(s, v)
}

func (s stack) pop() (stack, int) {
	if len(s) == 0 {
		log.Fatal("Popping empty stack!")
	}
	return s[:len(s)-1], s[len(s)-1]
}

// Interpreter represents a Brainf--k interpreter containing
// the program and memory which may be run.
type Interpreter struct {
	Program         []byte
	ProgramPosition int
	Memory          [memorySize]int64
	MemoryPosition  int
	bracketMap      map[int]int
	Output          strings.Builder
	Input           string
}

// Init initialises an interpreter by resetting the memory position.
func (ipr *Interpreter) Init() {
	ipr.MemoryPosition = (memorySize / 2)
}

// LoadProgram loads a program to the interpreter and
// builds a bracket map.
func (ipr *Interpreter) LoadProgram(data []byte) {
	ipr.bracketMap = make(map[int]int)
	var tempStack stack
	for i := 0; i < len(data); i++ {
		if data[i] == '>' || data[i] == '<' || data[i] == '+' ||
			data[i] == '-' || data[i] == '[' || data[i] == ']' || data[i] == '.' ||
			data[i] == ',' {
			ipr.Program = append(ipr.Program, data[i])
		}
		if data[i] == '[' {
			tempStack = tempStack.push(len(ipr.Program) - 1)
		} else if data[i] == ']' {
			var beginning int
			tempStack, beginning = tempStack.pop()
			ipr.bracketMap[beginning] = len(ipr.Program) - 1
			ipr.bracketMap[len(ipr.Program)-1] = beginning
		}
	}
	ipr.ProgramPosition = 0
}

// Clock runs one cycle/tick of the interpreter.
// It returns false when the program ends.
func (ipr *Interpreter) Clock() bool {
	switch ipr.Program[ipr.ProgramPosition] {
	case '>':
		ipr.MemoryPosition++
	case '<':
		ipr.MemoryPosition--
	case '+':
		ipr.Memory[ipr.MemoryPosition]++
	case '-':
		ipr.Memory[ipr.MemoryPosition]--
	case '.':
		fmt.Printf("%c", ipr.Memory[ipr.MemoryPosition])
		ipr.Output.WriteByte(byte(ipr.Memory[ipr.MemoryPosition]))
	case ',':
		if len(ipr.Input) > 0 {
			ipr.Memory[ipr.MemoryPosition] = int64(ipr.Input[0])
			ipr.Input = ipr.Input[1:]
		} else {
			fmt.Scanf("%c", &ipr.Memory[ipr.MemoryPosition])
		}
	case '[':
		if ipr.Memory[ipr.MemoryPosition] == 0 {
			ipr.ProgramPosition = ipr.bracketMap[ipr.ProgramPosition]
		}
	case ']':
		if ipr.Memory[ipr.MemoryPosition] != 0 {
			ipr.ProgramPosition = ipr.bracketMap[ipr.ProgramPosition]
		}
	}
	ipr.ProgramPosition++
	return !(ipr.ProgramPosition > len(ipr.Program)-1)
}
