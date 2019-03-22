package interpreter

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const memorySize int = 5000

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

// GetProperMemoryPosition returns the corrected
// memory position from the middle position of the memory array.
func (ipr *Interpreter) GetProperMemoryPosition() int {
	return (ipr.MemoryPosition - (memorySize / 2))
}

// GetProperMemoryValue returns a value from memory
// using a corrected memory position.
func (ipr *Interpreter) GetProperMemoryValue(pos int) int {
	return int(ipr.Memory[pos+(memorySize/2)])
}

// LoadFromFile loads a program to the interpreter from
// a file path.
func (ipr *Interpreter) LoadFromFile(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	ipr.LoadProgram(data)
}

// LoadProgram loads a program to the interpreter and
// builds a bracket map.
func (ipr *Interpreter) LoadProgram(data []byte) {
	ipr.MemoryPosition = (memorySize / 2)
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

// Run runs the program until the it ends.
func (ipr *Interpreter) Run() {
	for ipr.Clock() {
	}
}

// IsEnded returns whether the program has
// reached the end or not.
func (ipr *Interpreter) IsEnded() bool {
	return ipr.ProgramPosition > len(ipr.Program)-1
}

// Clock runs one cycle/tick of the interpreter.
// It returns false when the program ends.
func (ipr *Interpreter) Clock() bool {
	if ipr.IsEnded() {
		return false
	}
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
	return true
}
