package cartridge

type MBC0 struct {
	header *Header

	ROM []uint8
}

func (mbc *MBC0) RAMDump() []uint8 {
	return nil
}

func (mbc *MBC0) Header() *Header {
	return mbc.header
}

func NewMBC0(data []uint8, header *Header) *MBC0 {
	return &MBC0{
		header: header,
		ROM:    data,
	}
}

func (mbc *MBC0) Write(_ uint16, _ uint8) {
	// Writing to MBC0 has no effect
	return
}

func (mbc *MBC0) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		return mbc.ROM[addr]

	default:
		return 0xFF
	}
}
