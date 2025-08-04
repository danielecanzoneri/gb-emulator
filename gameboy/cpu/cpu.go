package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/mmu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
)

// Ticker describes hardware components that needs clock synchronization
type Ticker interface {
	Tick(ticks int)
}

type CPU struct {
	// 8-bit registers
	A, F uint8
	B, C uint8
	D, E uint8
	H, L uint8

	// 16-bit registers
	SP uint16
	PC uint16

	// Interrupt Master Enable
	IME        bool
	_EIDelayed bool // Set to true when EI is executed, but not yet effective

	// Flags to detect when interrupt is requested and cancelled (write to IE mid servicing)
	interruptMaskRequested       uint8
	writeIEHasCancelledInterrupt bool
	interruptCancelled           bool

	// Halt flag
	halted  bool
	haltBug bool

	// Other components
	mmu *mmu.MMU
	ppu *ppu.PPU

	cyclers []Ticker

	// Used for debugger
	callHook func()
	retHook  func()

	// Opcodes tables
	opcodesTable         [256]func()
	prefixedOpcodesTable [32]func(uint8)
}

func New(mmu *mmu.MMU, ppu *ppu.PPU) *CPU {
	cpu := &CPU{
		mmu: mmu,
		ppu: ppu,
	}
	cpu.initOpcodeTable()
	return cpu
}

func (cpu *CPU) AddCycler(cyclers ...Ticker) {
	for _, c := range cyclers {
		cpu.cyclers = append(cpu.cyclers, c)
	}
}

func (cpu *CPU) Tick(ticks int) {
	cpu.interruptCancelled = false

	// Cancel interrupt one cycle later
	if cpu.writeIEHasCancelledInterrupt {
		cpu.interruptCancelled = true
	}

	cpu.writeIEHasCancelledInterrupt = false
	if cpu.interruptMaskRequested > 0 && !(cpu.mmu.Read(ieAddr)&cpu.interruptMaskRequested > 0) {
		cpu.writeIEHasCancelledInterrupt = true
	}

	for _, cycler := range cpu.cyclers {
		cycler.Tick(ticks)
	}
}

func (cpu *CPU) ExecuteInstruction() {
	if !cpu.halted {
		opcode := cpu.ReadNextByte()

		// Execute opcode
		cpu.opcodesTable[opcode]()

	} else { // Cycle if halted
		cpu.Tick(4)
	}

	// Service eventual interrupts
	cpu.handleInterrupts()
}
