package main

import (
	"fmt"
	"log"
	"strings"
)

const memorySize int = 500

type stack []int

func (s stack) Push(v int) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, int) {
	if len(s) == 0 {
		log.Fatal("Popping empty stack!")
	}
	return s[:len(s)-1], s[len(s)-1]
}

type Interpreter struct {
	program         []byte
	programPosition int
	memory          [memorySize]int64
	memoryPosition  int
	bracketMap      map[int]int
	Output          strings.Builder
	Input           string
}

func (ipr *Interpreter) Init() {
	ipr.memoryPosition = (memorySize / 2)
}

func (ipr *Interpreter) LoadProgram(data []byte) {
	ipr.bracketMap = make(map[int]int)
	var tempStack stack
	for i := 0; i < len(data); i++ {
		if data[i] == '>' || data[i] == '<' || data[i] == '+' ||
			data[i] == '-' || data[i] == '[' || data[i] == ']' || data[i] == '.' ||
			data[i] == ',' {
			ipr.program = append(ipr.program, data[i])
		}
		if data[i] == '[' {
			tempStack = tempStack.Push(len(ipr.program) - 1)
		} else if data[i] == ']' {
			var beginning int
			tempStack, beginning = tempStack.Pop()
			ipr.bracketMap[beginning] = len(ipr.program) - 1
			ipr.bracketMap[len(ipr.program)-1] = beginning
		}
	}
	ipr.programPosition = 0
}

func (ipr *Interpreter) Clock() bool {
	switch ipr.program[ipr.programPosition] {
	case '>':
		ipr.memoryPosition++
	case '<':
		ipr.memoryPosition--
	case '+':
		ipr.memory[ipr.memoryPosition]++
	case '-':
		ipr.memory[ipr.memoryPosition]--
	case '.':
		fmt.Printf("%c", ipr.memory[ipr.memoryPosition])
		ipr.Output.WriteByte(byte(ipr.memory[ipr.memoryPosition]))
	case ',':
		if len(ipr.Input) > 0 {
			ipr.memory[ipr.memoryPosition] = int64(ipr.Input[0])
			ipr.Input = ipr.Input[1:]
		} else {
			fmt.Scanf("%c", &ipr.memory[ipr.memoryPosition])
		}
	case '[':
		if ipr.memory[ipr.memoryPosition] == 0 {
			ipr.programPosition = ipr.bracketMap[ipr.programPosition]
		}
	case ']':
		if ipr.memory[ipr.memoryPosition] != 0 {
			ipr.programPosition = ipr.bracketMap[ipr.programPosition]
		}
	}
	ipr.programPosition++
	return !(ipr.programPosition > len(ipr.program)-1)
}
