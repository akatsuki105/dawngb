package cpu

// P1/JOYP (0xFF00)
type Joypad struct {
	irq      func(n int)
	P14, P15 bool  // P1.4, P1.5 (P14, P15 は CPUのpinの名前)
	JOYP     uint8 // P1.0-3

	inputs uint8 // ゲームの実際のキー入力を反映したもの(pollで使用), (Dpad << 4) | Buttons, 0 is pressed, 1 is not pressed
}

func newJoypad(irq func(n int)) *Joypad {
	return &Joypad{
		irq:    irq,
		inputs: 0xFF,
	}
}

func (j *Joypad) reset() {
	j.P14, j.P15 = false, false
	j.JOYP = 0x0F
	j.inputs = 0xFF
}

// poll inputs
//
//	bit0-3: A, B, SELECT, START
//	bit4-7: RIGHT, LEFT, UP, DOWN
//
// 0 is pressed, 1 is not pressed
func (j *Joypad) poll(inputs uint8) {
	j.JOYP = 0x0F
	if !j.P14 {
		j.JOYP &= (inputs >> 4) & 0x0F
	}
	if !j.P15 {
		j.JOYP &= inputs & 0x0F
	}

	if j.JOYP != 0x0F {
		j.irq(IRQ_JOYPAD)
	}
}

func (j *Joypad) read() uint8 {
	j.poll(j.inputs)
	val := j.JOYP | 0xC0
	if j.P14 {
		val |= (1 << 4)
	}
	if j.P15 {
		val |= (1 << 5)
	}
	return val
}

func (j *Joypad) write(val uint8) {
	j.P14 = val&(1<<4) != 0
	j.P15 = val&(1<<5) != 0
}
