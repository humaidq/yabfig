package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

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

const memorySize int = 500

type Interpreter struct{
  program []byte
  programPosition int
  memory [memorySize]int64
  memoryPosition int
  bracketMap map[int]int
}

func (ipr *Interpreter) init() {
  ipr.memoryPosition = (memorySize/2)
}

func (ipr *Interpreter) loadProgram(data []byte) {
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

func (ipr *Interpreter) clock() bool {
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
  case ',':
    fmt.Scanf("%c", &ipr.memory[ipr.memoryPosition])
  case '[':
    if ipr.memory[ipr.memoryPosition] == 0 {
      ipr.programPosition = ipr.bracketMap[ipr.programPosition]
    }
  case ']':
    if ipr.memory[ipr.memoryPosition] != 0 {
      ipr.programPosition = ipr.bracketMap[ipr.programPosition]
    }
  }
  ipr.programPosition++;
  return !(ipr.programPosition > len(ipr.program)-2)
}

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
    ipr := Interpreter{}
    ipr.init()
    ipr.loadProgram(data)
    for ipr.clock() {
    }
	} else {
		fmt.Println("yabfig: usage: yabfig [file]")
	}
}
