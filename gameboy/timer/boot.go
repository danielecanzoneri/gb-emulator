package timer

func (t *Timer) SkipBoot() {
	t.systemCounter = 0xABC4
}
