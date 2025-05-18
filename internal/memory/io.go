package memory

const (
	JOYPAddr = 0xFF00

	DIVAddr  = 0xFF04
	TIMAAddr = 0xFF05
	TMAAddr  = 0xFF06
	TACAddr  = 0xFF07

	NR10Addr     = 0xFF10
	NR11Addr     = 0xFF11
	NR12Addr     = 0xFF12
	NR13Addr     = 0xFF13
	NR14Addr     = 0xFF14
	NR21Addr     = 0xFF16
	NR22Addr     = 0xFF17
	NR23Addr     = 0xFF18
	NR24Addr     = 0xFF19
	NR30Addr     = 0xFF1A
	NR31Addr     = 0xFF1B
	NR32Addr     = 0xFF1C
	NR33Addr     = 0xFF1D
	NR34Addr     = 0xFF1E
	NR41Addr     = 0xFF20
	NR42Addr     = 0xFF21
	NR43Addr     = 0xFF22
	NR44Addr     = 0xFF23
	NR50Addr     = 0xFF24
	NR51Addr     = 0xFF25
	NR52Addr     = 0xFF26
	WaveRAMAddr0 = 0xFF30
	WaveRAMAddr1 = 0xFF31
	WaveRAMAddr2 = 0xFF32
	WaveRAMAddr3 = 0xFF33
	WaveRAMAddr4 = 0xFF34
	WaveRAMAddr5 = 0xFF35
	WaveRAMAddr6 = 0xFF36
	WaveRAMAddr7 = 0xFF37
	WaveRAMAddr8 = 0xFF38
	WaveRAMAddr9 = 0xFF39
	WaveRAMAddrA = 0xFF3A
	WaveRAMAddrB = 0xFF3B
	WaveRAMAddrC = 0xFF3C
	WaveRAMAddrD = 0xFF3D
	WaveRAMAddrE = 0xFF3E
	WaveRAMAddrF = 0xFF3F

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
	case NR10Addr:
		fallthrough
	case NR11Addr:
		fallthrough
	case NR12Addr:
		fallthrough
	case NR13Addr:
		fallthrough
	case NR14Addr:
		fallthrough
	case NR21Addr:
		fallthrough
	case NR22Addr:
		fallthrough
	case NR23Addr:
		fallthrough
	case NR24Addr:
		fallthrough
	case NR30Addr:
		fallthrough
	case NR31Addr:
		fallthrough
	case NR32Addr:
		fallthrough
	case NR33Addr:
		fallthrough
	case NR34Addr:
		fallthrough
	case NR41Addr:
		fallthrough
	case NR42Addr:
		fallthrough
	case NR43Addr:
		fallthrough
	case NR44Addr:
		fallthrough
	case NR50Addr:
		fallthrough
	case NR51Addr:
		fallthrough
	case NR52Addr:
		fallthrough
	case WaveRAMAddr0:
		fallthrough
	case WaveRAMAddr1:
		fallthrough
	case WaveRAMAddr2:
		fallthrough
	case WaveRAMAddr3:
		fallthrough
	case WaveRAMAddr4:
		fallthrough
	case WaveRAMAddr5:
		fallthrough
	case WaveRAMAddr6:
		fallthrough
	case WaveRAMAddr7:
		fallthrough
	case WaveRAMAddr8:
		fallthrough
	case WaveRAMAddr9:
		fallthrough
	case WaveRAMAddrA:
		fallthrough
	case WaveRAMAddrB:
		fallthrough
	case WaveRAMAddrC:
		fallthrough
	case WaveRAMAddrD:
		fallthrough
	case WaveRAMAddrE:
		fallthrough
	case WaveRAMAddrF:
		mmu.APU.IOWrite(addr, v)

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

	// Interrupt flags
	case IFAddr:
		mmu.ifReg = v
	case IEAddr:
		mmu.ieReg = v

	default:
		mmu.Data[addr] = v
	}
}

func (mmu *MMU) readIO(addr uint16) uint8 {
	switch addr {
	// Joypad
	case JOYPAddr:
		return mmu.Joypad.Read()

	// Audio I/O
	case NR10Addr:
		fallthrough
	case NR11Addr:
		fallthrough
	case NR12Addr:
		fallthrough
	case NR13Addr:
		fallthrough
	case NR14Addr:
		fallthrough
	case NR21Addr:
		fallthrough
	case NR22Addr:
		fallthrough
	case NR23Addr:
		fallthrough
	case NR24Addr:
		fallthrough
	case NR30Addr:
		fallthrough
	case NR31Addr:
		fallthrough
	case NR32Addr:
		fallthrough
	case NR33Addr:
		fallthrough
	case NR34Addr:
		fallthrough
	case NR41Addr:
		fallthrough
	case NR42Addr:
		fallthrough
	case NR43Addr:
		fallthrough
	case NR44Addr:
		fallthrough
	case NR50Addr:
		fallthrough
	case NR51Addr:
		fallthrough
	case NR52Addr:
		fallthrough
	case WaveRAMAddr0:
		fallthrough
	case WaveRAMAddr1:
		fallthrough
	case WaveRAMAddr2:
		fallthrough
	case WaveRAMAddr3:
		fallthrough
	case WaveRAMAddr4:
		fallthrough
	case WaveRAMAddr5:
		fallthrough
	case WaveRAMAddr6:
		fallthrough
	case WaveRAMAddr7:
		fallthrough
	case WaveRAMAddr8:
		fallthrough
	case WaveRAMAddr9:
		fallthrough
	case WaveRAMAddrA:
		fallthrough
	case WaveRAMAddrB:
		fallthrough
	case WaveRAMAddrC:
		fallthrough
	case WaveRAMAddrD:
		fallthrough
	case WaveRAMAddrE:
		fallthrough
	case WaveRAMAddrF:
		return mmu.APU.IORead(addr)

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

	// DMA transfer
	case DMAAddr:
		return mmu.dmaReg

	// Interrupt flags
	case IFAddr:
		return ^interruptMask | (mmu.ifReg & interruptMask)
	case IEAddr:
		return ^interruptMask | (mmu.ieReg & interruptMask)

	default: // Unused I/O return 1
		return 0xFF
	}
}
