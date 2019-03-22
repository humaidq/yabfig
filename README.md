# Yet Another Brainf--k Interpreter

## 1. Description

yabfig is a [Brainf--k](https://en.wikipedia.org/wiki/Brainfuck) 
interpreter written in Go. It has also been extended to lint
code (by removing un-interpreted characters) and a gdb-style debugger.

## 2. Requirements

The following packages must be installed on your system.

- Go *(tested with 1.12)*
- Git

## 3. Copying and contributing

This program is written by Humaid AlQassimi,
and is distributed under the BSD 2 Clause license.  

Contribution is currently accepted as merge requests on GitLab. Patches
via email is also accepted if preferred.

## 4. Download and install

```sh
$ go get -u gitlab.com/humaid/yabfig
$ go install gitlab.com/humaid/yabfig
```

## 5. Usage

```sh
Usage: yabfig [option] <file>
Options:
	-lint		Lint (format) a Brainfuck file by removing spaces and non-instruction characters and output it to standard output.
	-debug		Run an interactive gdb-style debugger.
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

## 7. To-do

- [x] Unit testing & example programs
- [x] gdb-like debugger
  - [x] Viewing memory
  - [x] Stepping
  - [x] Breakpoints
  - [ ] Ignore
  - [ ] View source
- [ ] A web interface with debugger and visualiser
