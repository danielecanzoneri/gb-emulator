package timer

func (t *Timer) SkipDMGBoot() {
	t.systemCounter = 0xABCC
}

func (t *Timer) SkipCGBBoot() {
	// Since the boot ROM’s duration depends on the header’s contents
	// (and the player’s inputs in compatibility mode), this value is not reliable.
	t.systemCounter = 0x2678
}
