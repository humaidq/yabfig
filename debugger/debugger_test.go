package debugger

import (
	bf "gitlab.com/humaid/yabfig/interpreter"
	"testing"
)

func TestBreakpoint(t *testing.T) {
	program := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	dbg := Debugger{}
	dbg.interpreter = &bf.Interpreter{}
	dbg.interpreter.LoadProgram([]byte(program))
	dbg.init()
	dbg.breakpoints[94] = true
	dbg.runClocks()
	if dbg.interpreter.ProgramPosition != 94 {
		t.Errorf("Breakpoint is not the same as program position (%d vs %d)!", dbg.interpreter.ProgramPosition, 94)
	}
}

func TestWatchpoint(t *testing.T) {
	program := "++>>+<<<++++++++++++++++>>>"
	dbg := Debugger{}
	dbg.interpreter = &bf.Interpreter{}
	dbg.interpreter.LoadProgram([]byte(program))
	dbg.init()
	dbg.watchpoints[-1] = expression{greater, 5}
	dbg.runClocks()
	if dbg.interpreter.GetProperMemoryPosition() != -1 {
		t.Errorf("Watchpoint did not stop in correct position (%d vs %d)!", dbg.interpreter.MemoryPosition, -1)
	} else if dbg.interpreter.GetProperMemoryValue(-1) <= 5 {
		t.Errorf("Watchpoint did not stop at correct value (%d)!", dbg.interpreter.GetProperMemoryValue(-1))
	}

}
