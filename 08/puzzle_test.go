package main

import (
	"testing"
)

func TestInstructions(t *testing.T) {
	vm := &VM{
		PC:          0,
		Accumulator: 0,
	}

	insn := Instruction{
		Opcode: "acc",
		Arg:    1,
	}

	insn.Execute(vm)
	if vm.PC != 1 || vm.Accumulator != 1 {
		t.Errorf("%v: %#v", insn, *vm)
	}

	insn.Arg = -1
	vm.Reset()
	insn.Execute(vm)
	if vm.PC != 1 || vm.Accumulator != -1 {
		t.Errorf("%v: %#v", insn, *vm)
	}

	insn.Opcode = "jmp"
	insn.Arg = 10
	vm.Reset()
	insn.Execute(vm)
	if vm.PC != 10 || vm.Accumulator != 0 {
		t.Errorf("%v: %#v", insn, *vm)
	}
	insn.Arg = -10
	insn.Execute(vm)
	if vm.PC != 0 || vm.Accumulator != 0 {
		t.Errorf("%v: %#v", insn, *vm)
	}

	insn.Opcode = "nop"
	vm.Reset()
	insn.Execute(vm)
	if vm.PC != 1 || vm.Accumulator != 0 {
		t.Errorf("%v: %#v", insn, *vm)
	}
}

func TestProgram(t *testing.T) {
	vm := &VM{}

	program := Program{
		Instructions: []Instruction{
			Instruction{
				Opcode: "acc",
				Arg:    1,
			},
		},
	}

	vm.Execute(&program)
	if vm.PC != 1 || vm.Accumulator != 1 {
		t.Errorf("%#v", *vm)
	}

	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "acc",
		Arg:    -1,
	})

	vm.Reset()
	vm.Execute(&program)
	if vm.PC != 2 || vm.Accumulator != 0 {
		t.Errorf("%#v", *vm)
	}

	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "jmp",
		Arg:    4,
	})
	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "nop",
	})
	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "acc",
		Arg:    10,
	})
	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "nop",
	})
	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "nop",
	})
	program.Instructions = append(program.Instructions, Instruction{
		Opcode: "acc",
		Arg:    1,
	})

	vm.Reset()
	vm.Execute(&program)
	if vm.PC != 8 || vm.Accumulator != 1 {
		t.Errorf("%#v", *vm)
	}
}

func TestBreakpoint(t *testing.T) {
	vm := &VM{}

	program := Program{
		Instructions: []Instruction{
			Instruction{
				Opcode: "acc",
				Arg:    1,
			},
			Instruction{
				Opcode: "acc",
				Arg:    1,
			},
			Instruction{
				Opcode: "acc",
				Arg:    1,
			},
		},
	}

	vm.Breakpoint = func(vm *VM, insn *Instruction) bool {
		if vm.PC == 1 {
			insn.Arg = 10
		}
		return false
	}

	result := vm.Execute(&program)
	if !result {
		t.Error("expected normal termination")
	}
	if vm.PC != 3 || vm.Accumulator != 12 {
		t.Errorf("%#v", *vm)
	}

	vm.Breakpoint = func(vm *VM, insn *Instruction) bool {
		if vm.PC == 1 {
			return true
		}
		return false
	}

	vm.Reset()
	result = vm.Execute(&program)
	if result {
		t.Error("expected abnormal termination")
	}
	if vm.PC != 1 || vm.Accumulator != 1 {
		t.Errorf("%#v", *vm)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		s string
		i Instruction
	}{
		{
			s: "nop +0",
			i: Instruction{Opcode: "nop", Arg: 0},
		},
		{
			s: "acc +1",
			i: Instruction{Opcode: "acc", Arg: 1},
		},
		{
			s: "jmp -10",
			i: Instruction{Opcode: "jmp", Arg: -10},
		},
	}

	for _, test := range tests {
		i, err := ParseInstruction(test.s)
		if err != nil {
			t.Error(err)
		}

		if i != test.i {
			t.Errorf("expected %#v, got %#v", test.i, i)
		}
	}

}
