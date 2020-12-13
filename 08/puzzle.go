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
	return fmt.Sprintf("%s %d", insn.Opcode, insn.Arg)
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

type Tracer struct {
	Program *Program
	Indices []int
	Visited map[int]bool
}

func (t *Tracer) Reset() {
	t.Indices = t.Indices[:0]
	t.Visited = make(map[int]bool)
}

func (t *Tracer) Trace(vm *VM, insn *Instruction) bool {
	if t.Visited[vm.PC] {
		return true
	}

	t.Indices = append(t.Indices, vm.PC)
	t.Visited[vm.PC] = true
	return false
}

func (t *Tracer) Dump() {
	for _, i := range t.Indices {
		fmt.Println(t.Program.Instructions[i])
	}
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
	tracer := &Tracer{
		Program: program,
		Indices: make([]int, 0),
		Visited: make(map[int]bool),
	}
	vm.Breakpoint = func(vm *VM, insn *Instruction) bool { return tracer.Trace(vm, insn) }

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

	// Initial run should loop
	result := vm.Execute(program)
	if result {
		return fmt.Errorf("expected abnormal termination")
	}
	fmt.Println("Accumulator at breakpoint:", vm.Accumulator)

	lastIdx := tracer.Indices[len(tracer.Indices)-1]
	lastInsn := program.Instructions[lastIdx]
	fmt.Println("Last instruction before repeat:", lastIdx, lastInsn)

	if lastInsn.Opcode != "jmp" {
		return fmt.Errorf("last instruction not a jmp?")
	}

	// Patch all jumps to nops, and vice versa
	mapping := map[string]string{
		"nop": "jmp",
		"jmp": "nop",
	}
	for i := 0; i < len(program.Instructions); i++ {
		insn := program.Instructions[i]
		from := insn.Opcode
		if to, ok := mapping[insn.Opcode]; ok {
			insn.Opcode = to
			program.Instructions[i] = insn

			vm.Reset()
			tracer.Reset()
			result = vm.Execute(program)

			if result {
				fmt.Printf("Accumulator at normal termination (patched %s at %d): %d\n", from, i, vm.Accumulator)
				return nil
			}

			// Still looped, so revert
			insn.Opcode = from
			program.Instructions[i] = insn
		}
	}

	return fmt.Errorf("couldn't find a terminating case")
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
