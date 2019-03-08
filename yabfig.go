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

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
		var mem [100]int64
		var memPos int = (len(mem) / 2)
		init := func() ([]byte, map[int]int) {
			bracketmap := make(map[int]int)
			var tempStack stack
			var prog []byte
			for i := 0; i < len(data); i++ {
				if data[i] == '>' || data[i] == '<' || data[i] == '+' ||
					data[i] == '-' || data[i] == '[' || data[i] == ']' || data[i] == '.' ||
					data[i] == ',' {
					prog = append(prog, data[i])
				}
				if data[i] == '[' {
					tempStack = tempStack.Push(len(prog) - 1)
				} else if data[i] == ']' {
					var beginning int
					tempStack, beginning = tempStack.Pop()
					bracketmap[beginning] = len(prog) - 1
					bracketmap[len(prog)-1] = beginning
				}
			}
			return prog, bracketmap
		}
		prog, bracketmap := init()
		progPos := 0
		for ; ; progPos++ {
			switch prog[progPos] {
			case '>':
				memPos++
			case '<':
				memPos--
			case '+':
				mem[memPos]++
			case '-':
				mem[memPos]--
			case '.':
				fmt.Printf("%c", mem[memPos])
			case ',':
				fmt.Scanf("%c", &mem[memPos])
			case '[':
				if mem[memPos] == 0 {
					progPos = bracketmap[progPos]
				}
			case ']':
				if mem[memPos] != 0 {
					progPos = bracketmap[progPos]
				}
			}
			if progPos > len(prog)-2 {
				break
			}
		}
	} else {
		fmt.Println("yabfig: usage: yabfig [file]")
	}
}
