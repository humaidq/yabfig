# yabfig
## 1. Description
![Screenshot](https://humaidq.ae/projects/screenshots/yabfig.jpg)

yabfig is a [BF](https://en.wikipedia.org/wiki/brainfuck) 
interpreter written in Go. It has also been extended to lint
code (by removing un-interpreted characters) and to include a gdb-style
interpreter.

## 2. Requirements

The following packages must be installed on your system.

- Go *(tested with 1.12)*
- Git

## 3. Copying and contributing

This program is written by Humaid AlQassimi,
and is distributed under the
[BSD 2 Clause](https://humaidq.ae/license/bsd-2-clause) license.  

## 4. Download and install

```sh
$ go get -u git.sr.ht/~humaid/yabfig
$ go install git.sr.ht/~humaid/yabfig
```

## 5. Usage

```sh
Usage: yabfig [option] <file>
Options:
	-lint		Lint (format) a Brainfuck file by removing spaces and non-instruction characters and output it to standard output.
	-debug		Run an interactive gdb-style debugger.
```
To run the example program `hello-world.bf`:
```sh
$ yabfig programs/hello-world.bf
```
Using the debugger to set breakpoints:
```sh
$ yabfig -debug programs/hello-world.bf
yabfig debugger for Brainfuck.
Commands are similar to gdb, type "help" for a list of compatible commands.
(yabfig-dbg) help
List of commands:

run -- Run the program
print [pos] -- Print value at memory position
next [count] -- Execute next instruction[s]
jump [pos] -- Jump to a program position and resume
break [pos] -- Add breakpoint at program position
clear [pos] -- Delete breakpoint at program position
watch [n = x] -- Set watchpoint when memory position n is x
kill -- Kill program execution
(yabfig-dbg) b 98
Breakpoint #1 at position 98
(yabfig-dbg) b 102
Breakpoint #2 at position 102
(yabfig-dbg) b 106
Breakpoint #3 at position 106
(yabfig-dbg) r
Running program: programs/hello-world.bf
Hello WorldBreakpoint hit at position 98
(yabfig-dbg) c
!Breakpoint hit at position 102
(yabfig-dbg) c

Breakpoint hit at position 106
(yabfig-dbg) c
Program exited
(yabfig-dbg) 
```

## 6. Change log

- v0.1 *(Mar 8 2019)*
  - Initial release
- v0.2 *(Mar 18 2019)*
  - Added linter
  - Added unit tests
  - Interpreter as a struct with methods
- v0.2.1 *(Mar 18 2019)*
  - Add GoDoc
  - Move Interpreter to a separate package
- v0.3 *(Mar 22 2019)*
  - Add a simple gdb-style debugger
  - Improve interpreter functions

