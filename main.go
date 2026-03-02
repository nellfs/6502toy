package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nellfs/6502toy/internal/cpu"
)

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func main() {
	// assemble
	err := run("ca65", "test/minimal.asm", "-o", "build/minimal.o")
	if err != nil {
		panic(err)
	}

	err = run("ld65", "build/minimal.o", "-t", "none", "-o", "build/program.bin")
	if err != nil {
		panic(err)
	}

	fmt.Println("compiled to program.bin")

	cpu := &cpu.CPU{}
	cpu.PC = 0x8000
	data, err := os.ReadFile("build/program.bin")
	if err != nil {
		panic(err)
	}
	copy(cpu.Mem[0x8000:], data)
	cpu.Step()

	fmt.Printf("A = %02X\n", cpu.A)
}
