package cpu

import (
	"fmt"
	"log/slog"
)

type CPU struct {
	A   byte
	PC  uint16
	Mem [65536]byte
}

func (c *CPU) Step() error {
	opcode := c.Mem[c.PC]
	c.PC++

	switch opcode {
	case 0xA9:
		value := c.Mem[c.PC]
		c.PC++
		c.A = value
	default:
		slog.Error("step", "pc", fmt.Sprintf("0x%04X", c.PC), "opcode", fmt.Sprintf("0x%02X", opcode))
		return fmt.Errorf("unknown opcode: 0x%02X", opcode)
	}
	return nil
}

// Run executes instructions until a BRK (0x00) is encountered or
// the given maximum number of steps has been executed.
func (c *CPU) Run(maxSteps int) {
	for i := 0; i < maxSteps; i++ {
		opcode := c.Mem[c.PC]
		if opcode == 0x00 {
			slog.Info("halt (BRK)", "pc", c.PC)
			return
		}

		if err := c.Step(); err != nil {
			slog.Error("error in step", "error", err)
			return
		}
	}

	slog.Warn("warning: max steps reached", "maxSteps", maxSteps, "pc", c.PC)
}
