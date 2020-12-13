package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type VM struct {
	PC          int
	Accumulator int
	Breakpoint  func(*VM, *Instruction) bool
}

func (vm *VM) Reset() {
	vm.PC = 0
	vm.Accumulator = 0
}

func (vm *VM) Execute(program *Program) bool {
	normalTermination := true
	for vm.PC < len(program.Instructions) {
		insn := program.Instructions[vm.PC]
		if vm.Breakpoint != nil && vm.Breakpoint(vm, &insn) {
			normalTermination = false
			break
		}
		insn.Execute(vm)
	}

	return normalTermination
}

type Instruction struct {
	Opcode string
	Arg    int
}

func (insn Instruction) Execute(on *VM) {
	switch insn.Opcode {
	case "acc":
		on.Accumulator += insn.Arg
		on.PC += 1
	case "jmp":
		on.PC += insn.Arg
	case "nop":
		on.PC += 1
	default:
		panic(insn.Opcode)
	}
}

func (insn Instruction) String() string {
	switch insn.Opcode {
	case "acc", "jmp":
		return fmt.Sprintf("%s %d", insn.Opcode, insn.Arg)
	case "nop":
		return fmt.Sprintf("%s", insn.Opcode)
	default:
		panic(insn.Opcode)
	}

	return ""
}

func ParseInstruction(s string) (Instruction, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return Instruction{}, fmt.Errorf("couldn't parse instruction parts: %s", s)
	}

	arg, err := strconv.Atoi(parts[1])
	if err != nil {
		return Instruction{}, err
	}

	return Instruction{
		Opcode: parts[0],
		Arg:    arg,
	}, nil
}

type Program struct {
	Instructions []Instruction
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: %s INPUT", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	vm := &VM{}
	program := &Program{
		Instructions: make([]Instruction, 0, 100),
	}

	visited := make(map[int]bool)
	bp := func(vm *VM, insn *Instruction) bool {
		if visited[vm.PC] {
			return true
		}
		visited[vm.PC] = true
		return false
	}
	vm.Breakpoint = bp

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		i, err := ParseInstruction(line)
		if err != nil {
			return err
		}

		program.Instructions = append(program.Instructions, i)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("Parsed", len(program.Instructions))

	result := vm.Execute(program)
	if result {
		return fmt.Errorf("expected abnormal termination")
	}
	fmt.Println("Accumulator at breakpoint:", vm.Accumulator)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
