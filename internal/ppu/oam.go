package ppu

const (
	OAMSize = 0xA0
)

type OAM struct {
	data [OAMSize]uint8
}

func (o *OAM) read(addr uint16) uint8 {
	return o.data[addr]
}

func (o *OAM) write(addr uint16, value uint8) {
	o.data[addr] = value
}
