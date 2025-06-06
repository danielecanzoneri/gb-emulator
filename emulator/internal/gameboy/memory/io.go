package memory

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

	IFAddr              = 0xFF0F
	IEAddr              = 0xFFFF
	interruptMask uint8 = 0b00011111
)

func (mmu *MMU) writeIO(addr uint16, v uint8) {
	switch addr {
	// Joypad
	case JOYPAddr:
		mmu.Joypad.Write(v)

	// Audio I/O
	case NR10Addr, NR11Addr, NR12Addr, NR13Addr, NR14Addr,
		NR21Addr, NR22Addr, NR23Addr, NR24Addr,
		NR30Addr, NR31Addr, NR32Addr, NR33Addr, NR34Addr,
		NR41Addr, NR42Addr, NR43Addr, NR44Addr,
		NR50Addr, NR51Addr, NR52Addr:
		mmu.APU.IOWrite(addr, v)

	// Timer I/O
	case DIVAddr, TIMAAddr, TMAAddr, TACAddr:
		mmu.Timer.Write(addr, v)

	// PPU I/O
	case LCDCAddr, STATAddr, SCYAddr, SCXAddr, LYAddr, LYCAddr, BGPAddr, OBP0Addr, OBP1Addr, WYAddr, WXAddr:
		mmu.PPU.Write(addr, v)

	// DMA transfer
	case DMAAddr:
		mmu.DMA(v)

	// Interrupt flags
	case IFAddr:
		mmu.ifReg = v
	case IEAddr:
		mmu.ieReg = v

	default:
		// Wave RAM
		if waveRAMStartAddr <= addr && addr < waveRAMStartAddr+waveRAMLength {
			mmu.APU.IOWrite(addr, v)
		}
	}
}

func (mmu *MMU) readIO(addr uint16) uint8 {
	switch addr {
	// Joypad
	case JOYPAddr:
		return mmu.Joypad.Read()

	// Audio I/O
	case NR10Addr, NR11Addr, NR12Addr, NR13Addr, NR14Addr,
		NR21Addr, NR22Addr, NR23Addr, NR24Addr,
		NR30Addr, NR31Addr, NR32Addr, NR33Addr, NR34Addr,
		NR41Addr, NR42Addr, NR43Addr, NR44Addr,
		NR50Addr, NR51Addr, NR52Addr:
		return mmu.APU.IORead(addr)

	// Timer I/O
	case DIVAddr, TIMAAddr, TMAAddr, TACAddr:
		return mmu.Timer.Read(addr)

	// PPU I/O
	case LCDCAddr, STATAddr, SCYAddr, SCXAddr, LYAddr, LYCAddr, BGPAddr, OBP0Addr, OBP1Addr, WYAddr, WXAddr:
		return mmu.PPU.Read(addr)

	// DMA transfer
	case DMAAddr:
		return mmu.dmaReg

	// Interrupt flags
	case IFAddr:
		return ^interruptMask | (mmu.ifReg & interruptMask)
	case IEAddr:
		return ^interruptMask | (mmu.ieReg & interruptMask)

	default:
		// Wave RAM
		if waveRAMStartAddr <= addr && addr < waveRAMStartAddr+waveRAMLength {
			return mmu.APU.IORead(addr)
		}

		// Unused I/O return bits 1
		return 0xFF
	}
}
