package audio

import "github.com/akatsuki105/dugb/util"

func (a *audio) Read(addr uint16) uint8 {
	a.CatchUp()
	return a.ioreg[addr-0xFF10]
}

func (a *audio) Write(addr uint16, val uint8) {
	a.CatchUp()

	a.ioreg[addr-0xFF10] = val
	switch addr {
	case 0xFF10:
		a.ch1.sweep.speed = int(val>>4) & 0b111
		if a.ch1.sweep.speed == 0 {
			a.ch1.sweep.speed = 8
		}
		a.ch1.sweep.step = a.ch1.sweep.speed
		a.ch1.sweep.up = !util.Bit(val, 3)
		a.ch1.sweep.shift = int(val & 0b111)

	case 0xFF11:
		a.ch1.duty = int(val >> 6)
		a.ch1.length = 64 - int(val&0b11_1111)

	case 0xFF12:
		a.ch1.envelope.initialVolume = int(val>>4) & 0b1111
		a.ch1.envelope.direction = util.Bit(val, 3)
		a.ch1.envelope.speed = int(val & 0b111)

	case 0xFF13:
		a.ch1.period = (a.ch1.period & 0xFF00) | int(val)
		a.ch1.freqCounter = a.ch1.dutyStepCycle()

	case 0xFF14:
		a.ch1.period = (a.ch1.period & 0x00FF) | (int(val&0b111) << 8)
		a.ch1.freqCounter = a.ch1.dutyStepCycle()
		a.ch1.stop = util.Bit(val, 6)
		if util.Bit(val, 7) {
			a.ch1.enabled = true
			a.ch1.envelope.reset()
			a.ch1.sweep.reset()
			if a.ch1.length == 0 {
				a.ch1.length = 64
			}
		}

	case 0xFF16:
		a.ch2.duty = int(val >> 6)
		a.ch2.length = 64 - int(val&0b11_1111)

	case 0xFF17:
		a.ch2.envelope.initialVolume = int(val>>4) & 0b1111
		a.ch2.envelope.direction = util.Bit(val, 3)
		a.ch2.envelope.speed = int(val & 0b111)

	case 0xFF18:
		a.ch2.period = (a.ch2.period & 0xFF00) | int(val)
		a.ch2.freqCounter = a.ch2.dutyStepCycle()

	case 0xFF19:
		a.ch2.period = (a.ch2.period & 0x00FF) | (int(val&0b111) << 8)
		a.ch2.freqCounter = a.ch2.dutyStepCycle()
		a.ch2.stop = util.Bit(val, 6)
		if util.Bit(val, 7) {
			a.ch2.enabled = true
			a.ch2.envelope.reset()
			if a.ch2.length == 0 {
				a.ch2.length = 64
			}
		}

	case 0xFF1A:
		a.ch3.active = util.Bit(val, 7)
	case 0xFF1B:
		a.ch3.length = 256 - int(val)
	case 0xFF1C:
		a.ch3.volume = [4]int{4, 0, 1, 2}[int(val>>5)&0b11] // 波形は最大15なので4左シフトすれば0%
	case 0xFF1D:
		a.ch3.period &= 0b111_0000_0000
		a.ch3.period |= int(val)
		a.ch3.freqCounter = a.ch3.windowStepCycle()
	case 0xFF1E:
		a.ch3.stop = util.Bit(val, 6)
		a.ch3.period &= 0b000_1111_1111
		a.ch3.period |= int(val&0b111) << 8
		a.ch3.freqCounter = a.ch3.windowStepCycle()
		if util.Bit(val, 7) {
			a.ch3.enabled = a.ch3.active
			if a.ch3.length == 0 {
				a.ch3.length = 256
			}
			a.ch3.window = 0
		}

	case 0xFF26:
		a.enabled = util.Bit(val, 7)

	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		a.ch3.samples[addr-0xFF30] = val
	}
}