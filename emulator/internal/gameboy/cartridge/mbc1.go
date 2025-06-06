package cartridge

import "log"

type MBC1 struct {
	header *Header

	ROMBanks uint
	RAMBanks uint

	ROM []uint8
	RAM []uint8

	// Registers
	RAMEnabled     bool
	currentROMBank uint  // 5 bit register
	currentRAMBank uint  // 2 bit register (can be used as high bits of currentROMBank if number of banks is > 32)
	mode           uint8 // 0: ROM and RAM banks 0 fixed, 1: ROM bank 0 and RAM banks can be switched

	useRAMRegisterAsHighROMRegister bool
}

func (mbc *MBC1) Header() *Header {
	return mbc.header
}

func NewMBC1(data []uint8, header *Header) *MBC1 {
	mbc := &MBC1{
		header:                          header,
		ROMBanks:                        header.ROMBanks,
		RAMBanks:                        header.RAMBanks,
		ROM:                             data,
		currentROMBank:                  1,
		useRAMRegisterAsHighROMRegister: header.ROMBanks > 32,
	}
	// Always have a RAM bank
	if mbc.RAMBanks == 0 {
		mbc.RAMBanks = 1
	}
	mbc.RAM = make([]uint8, mbc.RAMBanks*0x2000)

	return mbc
}

func (mbc *MBC1) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		mbc.enableRAM(value)
	case addr < 0x4000:
		mbc.SetROMBank(value)
	case addr < 0x6000:
		mbc.SetRAMBank(value)
	case addr < 0x8000:
		mbc.SetMode(value)

	case 0xA000 <= addr && addr < 0xC000:
		RAMAddress := mbc.computeRAMAddress(addr)
		mbc.RAM[RAMAddress] = value

	default:
		log.Printf("[WARN] MBC1 Write address is out of range: %04X\n", addr)
	}
}

func (mbc *MBC1) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := mbc.computeROMAddress(addr)
		return mbc.ROM[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		RAMAddress := mbc.computeRAMAddress(addr)
		return mbc.RAM[RAMAddress]

	default:
		log.Printf("[WARN] Cartridge Read address is out of range: %04X\n", addr)
		return 0xFF
	}
}

func (mbc *MBC1) enableRAM(value uint8) {
	// Low nibble = 0xA enables RAM
	mbc.RAMEnabled = value&0x0F == 0xA
}

func (mbc *MBC1) SetROMBank(value uint8) {
	// Only lower 5 bits are used
	value = value & 0x1F

	// Register 0 behaves as 1
	if value == 0 {
		value = 1
	}

	// Number of bank is masked to the required number of bits
	mbc.currentROMBank = uint(value) % mbc.ROMBanks
}

func (mbc *MBC1) SetRAMBank(value uint8) {
	// Only lower 2 bits are used
	value = value & 0x3

	// Set only if ROM or RAM are large enough
	if mbc.RAMBanks > 1 || mbc.useRAMRegisterAsHighROMRegister {
		mbc.currentRAMBank = uint(value)
	}
}

func (mbc *MBC1) SetMode(value uint8) {
	// 1 bit register
	mbc.mode = value & 1
}

func (mbc *MBC1) computeROMAddress(cpuAddress uint16) uint {
	// bank number: 2 bits - 5 bits, cpuAddress: 14 bits
	var bankNumber uint = 0
	switch {
	case cpuAddress < 0x4000:
		if mbc.mode == 1 && mbc.useRAMRegisterAsHighROMRegister {
			bankNumber = mbc.currentRAMBank << 5
		}
	case cpuAddress < 0x8000:
		// Select 14 bits
		cpuAddress = cpuAddress & 0x3FFF

		bankNumber = mbc.currentROMBank
		if mbc.useRAMRegisterAsHighROMRegister {
			bankNumber |= mbc.currentRAMBank << 5
		}
	default:
		panic("should never happen")
	}
	return bankNumber<<14 | uint(cpuAddress)
}

func (mbc *MBC1) computeRAMAddress(cpuAddress uint16) uint {
	// bank number: 2 bits, cpuAddress: 13 bits
	switch {
	case 0xA000 <= cpuAddress && cpuAddress < 0xC000:
		cpuAddress = cpuAddress & 0x1FFF
		if !mbc.useRAMRegisterAsHighROMRegister {
			return mbc.currentRAMBank<<13 | uint(cpuAddress)
		}
		return uint(cpuAddress)
	default:
		panic("should never happen")
	}
}
