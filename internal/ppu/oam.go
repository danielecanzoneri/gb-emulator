package ppu

const (
	OAMSize = 0xA0
)

type OAM struct {
	Data [OAMSize]uint8
}

func (o *OAM) read(addr uint16) uint8 {
	return o.Data[addr]
}

func (o *OAM) write(addr uint16, value uint8) {
	o.Data[addr] = value
}
