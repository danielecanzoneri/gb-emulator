package mmu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/serial"
)

const (
	JOYPAddr = 0xFF00

	DIVAddr  = 0xFF04
	TIMAAddr = 0xFF05
	TMAAddr  = 0xFF06
	TACAddr  = 0xFF07

	NR10Addr         = 0xFF10
	NR11Addr         = 0xFF11
	NR12Addr         = 0xFF12
	NR13Addr         = 0xFF13
	NR14Addr         = 0xFF14
	NR21Addr         = 0xFF16
	NR22Addr         = 0xFF17
	NR23Addr         = 0xFF18
	NR24Addr         = 0xFF19
	NR30Addr         = 0xFF1A
	NR31Addr         = 0xFF1B
	NR32Addr         = 0xFF1C
	NR33Addr         = 0xFF1D
	NR34Addr         = 0xFF1E
	NR41Addr         = 0xFF20
	NR42Addr         = 0xFF21
	NR43Addr         = 0xFF22
	NR44Addr         = 0xFF23
	NR50Addr         = 0xFF24
	NR51Addr         = 0xFF25
	NR52Addr         = 0xFF26
	waveRAMStartAddr = 0xFF30
	waveRAMLength    = 16

	LCDCAddr = 0xFF40
	STATAddr = 0xFF41
	SCYAddr  = 0xFF42
	SCXAddr  = 0xFF43
	LYAddr   = 0xFF44
	LYCAddr  = 0xFF45
	DMAAddr  = 0xFF46
	BGPAddr  = 0xFF47
	OBP0Addr = 0xFF48
	OBP1Addr = 0xFF49
	WYAddr   = 0xFF4A
	WXAddr   = 0xFF4B
	VBKAddr  = 0xFF4F

	BANKAddr = 0xFF50

	HDMA1Addr = 0xFF51
	HDMA2Addr = 0xFF52
	HDMA3Addr = 0xFF53
	HDMA4Addr = 0xFF54
	HDMA5Addr = 0xFF55

	BGPIAddr = 0xFF68
	BGPDAddr = 0xFF69
	OBPIAddr = 0xFF6A
	OBPDAddr = 0xFF6B

	WBKAddr = 0xFF70

	IFAddr       = 0xFF0F
	IEAddr       = 0xFFFF
	ifMask uint8 = 0b00011111
)

func (mmu *MMU) writeIO(addr uint16, v uint8) {
	switch addr {
	// Joypad
	case JOYPAddr:
		mmu.joypad.Write(v)

	// Serial
	case serial.SBAddr, serial.SCAddr:
		mmu.serial.Write(addr, v)

	// Audio I/O
	case NR10Addr, NR11Addr, NR12Addr, NR13Addr, NR14Addr,
		NR21Addr, NR22Addr, NR23Addr, NR24Addr,
		NR30Addr, NR31Addr, NR32Addr, NR33Addr, NR34Addr,
		NR41Addr, NR42Addr, NR43Addr, NR44Addr,
		NR50Addr, NR51Addr, NR52Addr:
		mmu.apu.IOWrite(addr, v)

	// Timer I/O
	case DIVAddr, TIMAAddr, TMAAddr, TACAddr:
		mmu.timer.Write(addr, v)

	// PPU I/O
	case LCDCAddr, STATAddr, SCYAddr, SCXAddr, LYAddr, LYCAddr, BGPAddr, OBP0Addr, OBP1Addr, WYAddr, WXAddr, VBKAddr,
		BGPIAddr, BGPDAddr, OBPIAddr, OBPDAddr:
		mmu.ppu.Write(addr, v)

	// DMA transfer
	case DMAAddr:
		mmu.DMA(v)

	// Disable BOOT ROM
	case BANKAddr:
		mmu.DisableBootROM()

	// VDMA transfer
	case HDMA1Addr:
		mmu.vDMASrcHigh = v
	case HDMA2Addr:
		mmu.vDMASrcLow = v & 0xF0
	case HDMA3Addr:
		mmu.vDMADestHigh = v & 0x1F
	case HDMA4Addr:
		mmu.vDMADestLow = v & 0xF0
	case HDMA5Addr:
		mmu.VDMA(v)

	// wRAM bank register
	case WBKAddr:
		mmu.vbk = v & 0b111

	// Interrupt flags
	case IFAddr:
		mmu.ifReg = v
	case IEAddr:
		mmu.ieReg = v

	default:
		// Wave RAM
		if waveRAMStartAddr <= addr && addr < waveRAMStartAddr+waveRAMLength {
			mmu.apu.IOWrite(addr, v)
		}
	}
}

func (mmu *MMU) readIO(addr uint16) uint8 {
	switch addr {
	// Joypad
	case JOYPAddr:
		return mmu.joypad.Read()

	// Serial
	case serial.SBAddr, serial.SCAddr:
		return mmu.serial.Read(addr)

	// Audio I/O
	case NR10Addr, NR11Addr, NR12Addr, NR13Addr, NR14Addr,
		NR21Addr, NR22Addr, NR23Addr, NR24Addr,
		NR30Addr, NR31Addr, NR32Addr, NR33Addr, NR34Addr,
		NR41Addr, NR42Addr, NR43Addr, NR44Addr,
		NR50Addr, NR51Addr, NR52Addr:
		return mmu.apu.IORead(addr)

	// Timer I/O
	case DIVAddr, TIMAAddr, TMAAddr, TACAddr:
		return mmu.timer.Read(addr)

	// PPU I/O
	case LCDCAddr, STATAddr, SCYAddr, SCXAddr, LYAddr, LYCAddr, BGPAddr, OBP0Addr, OBP1Addr, WYAddr, WXAddr, VBKAddr,
		BGPIAddr, BGPDAddr, OBPIAddr, OBPDAddr:
		return mmu.ppu.Read(addr)

	// DMA transfer
	case DMAAddr:
		return mmu.dmaReg

	// VDMA transfer
	case HDMA5Addr:
		if mmu.vDMAActive {
			return mmu.vDMALength
		} else {
			return 0x80 | mmu.vDMALength
		}

	// wRAM bank register
	case WBKAddr:
		return 0xF8 | mmu.vbk

	// Interrupt flags
	case IFAddr:
		return ^ifMask | (mmu.ifReg & ifMask)
	case IEAddr:
		return mmu.ieReg

	default:
		// Wave RAM
		if waveRAMStartAddr <= addr && addr < waveRAMStartAddr+waveRAMLength {
			return mmu.apu.IORead(addr)
		}

		// Unused I/O return bits 1
		return 0xFF
	}
}
