package memory

import "fmt"

const (
	vRAM           = 0x8000
	eRAM           = 0xA000
	echoRAM        = 0xE000
	OAM            = 0xFE00
	reservedMemory = 0xFEA0
	ioRegisters    = 0xFF00
	hRAM           = 0xFF80
)

const (
	DIVAddr  = 0xFF04
	TIMAAddr = 0xFF05
	TMAAddr  = 0xFF06
	TACAddr  = 0xFF07

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
)

func (mmu *MMU) writeIO(addr uint16, v uint8) {
	switch addr {
	// Timer I/O
	case DIVAddr:
		fallthrough
	case TIMAAddr:
		fallthrough
	case TMAAddr:
		fallthrough
	case TACAddr:
		mmu.Timer.Write(addr, v)

	// PPU I/O
	case LCDCAddr:
		fallthrough
	case STATAddr:
		fallthrough
	case SCYAddr:
		fallthrough
	case SCXAddr:
		fallthrough
	case LYAddr:
		fallthrough
	case LYCAddr:
		fallthrough
	case BGPAddr:
		fallthrough
	case OBP0Addr:
		fallthrough
	case OBP1Addr:
		fallthrough
	case WYAddr:
		fallthrough
	case WXAddr:
		mmu.PPU.Write(addr, v)

	// DMA transfer
	case DMAAddr:
		mmu.DMA(v)

	default:
		mmu.Data[addr] = v
	}

	// Debug on serial port
	if addr == 0xFF01 {
		if v == 0 {
			fmt.Println()
		} else {
			fmt.Printf("%c", v)
		}
	}
}

func (mmu *MMU) readIO(addr uint16) uint8 {
	switch addr {
	// Timer I/O
	case DIVAddr:
		fallthrough
	case TIMAAddr:
		fallthrough
	case TMAAddr:
		fallthrough
	case TACAddr:
		return mmu.Timer.Read(addr)

	// PPU I/O
	case LCDCAddr:
		fallthrough
	case STATAddr:
		fallthrough
	case SCYAddr:
		fallthrough
	case SCXAddr:
		fallthrough
	case LYAddr:
		fallthrough
	case LYCAddr:
		fallthrough
	case BGPAddr:
		fallthrough
	case OBP0Addr:
		fallthrough
	case OBP1Addr:
		fallthrough
	case WYAddr:
		fallthrough
	case WXAddr:
		return mmu.PPU.Read(addr)

	default:
		return mmu.Data[addr]
	}
}
