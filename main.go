package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/nellfs/6502toy/internal/cpu"
)

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func main() {
	// window setup
	app := adw.NewApplication("com.nellfs.6502toy", gio.ApplicationFlagsNone)

	app.ConnectActivate(func() { activate(app) })
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}

	// code compilation

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

	// cpu runtime

	cpu := &cpu.CPU{}
	cpu.PC = 0x8000
	data, err := os.ReadFile("build/program.bin")
	if err != nil {
		panic(err)
	}
	copy(cpu.Mem[0x8000:], data)
	cpu.Step()

	fmt.Printf("accumulator = %02X\n", cpu.A)
}

func activate(app *adw.Application) {
	win := adw.NewApplicationWindow(&app.Application)
	win.SetTitle("6502 Toy")
	win.SetSizeRequest(320, 320)
	win.SetDefaultSize(800, 600)
	win.Present()
}
