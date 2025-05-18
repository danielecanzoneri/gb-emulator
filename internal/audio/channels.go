package audio

type Channel interface {
	IsActive() bool
	Output() float32

	Cycle()
	Disable()

	ReadRegister(addr uint16) uint8
	WriteRegister(addr uint16, value uint8)
}

var waveforms = [4][8]bool{
	//                       waveform                         | wave duty | duty cycle
	{false, false, false, false, false, false, false, true}, //        00 |     12.5 %
	{true, false, false, false, false, false, false, true},  //        01 |       25 %
	{true, false, false, false, false, true, true, true},    //        10 |       50 %
	{false, true, true, true, true, true, true, false},      //        11 |       75 %
}

var waveformsU8 = [4][8]uint8{
	//         waveform        | wave duty | duty cycle
	{0, 0, 0, 0, 0, 0, 0, 1}, //        00 |     12.5 %
	{1, 0, 0, 0, 0, 0, 0, 1}, //        01 |       25 %
	{1, 0, 0, 0, 0, 1, 1, 1}, //        10 |       50 %
	{0, 1, 1, 1, 1, 1, 1, 0}, //        11 |       75 %
}
