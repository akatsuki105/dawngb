package cpu

// P1/JOYP (0xFF00)
type joypad struct {
	irq      func(n int)
	p14, p15 bool  // P1.4, P1.5 (P14, P15 は CPUのpinの名前)
	joyp     uint8 // P1.0-3

	inputs uint8 // ゲームの実際のキー入力を反映したもの(pollで使用), (Dpad << 4) | Buttons, 0 is pressed, 1 is not pressed
}

func newJoypad(irq func(n int)) *joypad {
	return &joypad{
		irq:    irq,
		inputs: 0xFF,
	}
}

func (j *joypad) reset(hasBIOS bool) {
	j.p14, j.p15 = false, false
	j.joyp = 0x0F
	j.inputs = 0xFF
	if !hasBIOS {
		j.write(0x30)
		j.write(0xCF)
	}
}

// poll inputs
//
//	bit0-3: A, B, SELECT, START
//	bit4-7: RIGHT, LEFT, UP, DOWN
//
// 0 is pressed, 1 is not pressed
func (j *joypad) poll(inputs uint8) {
	j.joyp = 0x0F
	if !j.p14 {
		j.joyp &= (inputs >> 4) & 0x0F
	}
	if !j.p15 {
		j.joyp &= inputs & 0x0F
	}

	if j.joyp != 0x0F {
		j.irq(IRQ_JOYPAD)
	}
}

func (j *joypad) read() uint8 {
	j.poll(j.inputs)
	val := j.joyp | 0xC0
	if j.p14 {
		val |= (1 << 4)
	}
	if j.p15 {
		val |= (1 << 5)
	}
	return val
}

func (j *joypad) write(val uint8) {
	j.p14 = val&(1<<4) != 0
	j.p15 = val&(1<<5) != 0
}
