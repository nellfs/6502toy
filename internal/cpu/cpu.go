package cpu

import "log/slog"

type CPU struct {
	A   byte
	PC  uint16
	Mem [65536]byte
}

func (c *CPU) Step() {
	opcode := c.Mem[c.PC]
	c.PC++

	switch opcode {
	case 0xA9:
		value := c.Mem[c.PC]
		c.PC++
		c.A = value
	default:
		slog.Error("unknown opcode", "opcode", opcode, "pc", c.PC-1)
	}
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

		c.Step()
	}

	slog.Warn("warning: max steps reached", "maxSteps", maxSteps, "pc", c.PC)
}
