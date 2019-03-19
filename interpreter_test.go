package main

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	program := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	ipr := Interpreter{}
	ipr.Init()
	ipr.LoadProgram([]byte(program))
	for ipr.Clock() {
	}
	if ipr.Output.String() != "Hello World!\n" {
		t.Errorf("Hello world program output unexpected: \"%s\"", ipr.Output.String())
	}

}

func TestRot13(t *testing.T) {
	program := "-,+[-[>>++++[>++++++++<-]<+<-[>+>+>-[>>>]<[[>+<-]>>+>]<<<<<-]]>>>[-]+>--[-[<->+++[-]]]<[++++++++++++<[>-[>+>>]>[+[<+>-]>+>>]<<<<<-]>>[<+>-]>[-[-<<[-]>>]<<[<<->>-]>>]<<[<<+>>-]]<[-]<.[-]<-,+]"
	ipr := Interpreter{}
	ipr.Init()
	ipr.LoadProgram([]byte(program))
	ipr.Input = "Hello rot13! "
	for ipr.Clock() {
		if len(ipr.Input) == 0 {
			break
		}
	}
	if ipr.Output.String() != "Uryyb ebg13!" {
		t.Errorf("Rot13 program output unexpected: \"%s\"", ipr.Output.String())
	}

}
