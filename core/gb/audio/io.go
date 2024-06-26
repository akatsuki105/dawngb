package audio

import (
	"github.com/akatsuki105/dawngb/util"
)

// NOTE: "ゼルダの伝説 ふしぎの木の実" はAPUのNR52を正しく実装しないとタイトル画面から進めない

func (a *audio) Read(addr uint16) uint8 {
	a.CatchUp()
	switch addr {
	case 0xFF26:
		val := uint8(0)
		val = util.SetBit(val, 7, a.enabled)
		val = util.SetBit(val, 0, a.ch1.enabled)
		val = util.SetBit(val, 1, a.ch2.enabled)
		val = util.SetBit(val, 2, a.ch3.enabled)
		val = util.SetBit(val, 3, a.ch4.enabled)
		return val
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		if a.model == APU_GBA {
			if a.ch3.bank == 0 {
				return a.ch3.samples[16+addr-0xFF30]
			} else {
				return a.ch3.samples[addr-0xFF30]
			}
		}
	}
	return a.ioreg[addr-0xFF10]
}

func (a *audio) Write(addr uint16, val uint8) {
	a.CatchUp()

	if addr == 0xFF26 {
		a.enabled = util.Bit(val, 7)
	}

	if !a.enabled {
		return
	}

	a.ioreg[addr-0xFF10] = val
	switch addr {
	case 0xFF10:
		negate := util.Bit(val, 3)
		if a.ch1.sweep.enabled && a.ch1.sweep.negate && !negate {
			a.ch1.enabled = false // 下降スイープ中に上昇スイープを設定すると消音される
		}
		a.ch1.sweep.shift = int(val & 0b111)
		a.ch1.sweep.negate = negate
		a.ch1.sweep.interval = int(val>>4) & 0b111

	case 0xFF11:
		a.ch1.duty = int(val >> 6)
		a.ch1.length = 64 - int(val&0b11_1111)

	case 0xFF12:
		a.ch1.envelope.initialVolume = int(val>>4) & 0b1111
		a.ch1.envelope.direction = util.Bit(val, 3)
		a.ch1.envelope.speed = int(val & 0b111)
		if !a.ch1.dacEnable() {
			a.ch1.enabled = false
		}

	case 0xFF13:
		a.ch1.period = (a.ch1.period & 0xFF00) | int(val)

	case 0xFF14:
		a.ch1.period = (a.ch1.period & 0x00FF) | (int(val&0b111) << 8)
		a.ch1.stop = util.Bit(val, 6)
		if util.Bit(val, 7) { // キーオン(音が鳴り始める)
			a.ch1.tryRestart()
		}

	case 0xFF16:
		a.ch2.duty = int(val >> 6)
		a.ch2.length = 64 - int(val&0b11_1111)

	case 0xFF17:
		a.ch2.envelope.initialVolume = int(val>>4) & 0b1111
		a.ch2.envelope.direction = util.Bit(val, 3)
		a.ch2.envelope.speed = int(val & 0b111)
		if !a.ch2.dacEnable() {
			a.ch2.enabled = false
		}

	case 0xFF18:
		a.ch2.period = (a.ch2.period & 0xFF00) | int(val)

	case 0xFF19:
		a.ch2.period = (a.ch2.period & 0x00FF) | (int(val&0b111) << 8)
		a.ch2.stop = util.Bit(val, 6)
		if util.Bit(val, 7) { // キーオン(音が鳴り始める)
			a.ch2.tryRestart()
		}

	case 0xFF1A:
		a.ch3.dacEnable = util.Bit(val, 7)
		if !a.ch3.dacEnable {
			a.ch3.enabled = false
		}
		if a.model == APU_GBA {
			a.ch3.mode = int(val>>5) & 0b1
			a.ch3.bank = int(val>>6) & 0b1
		}
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
		if util.Bit(val, 7) { // キーオン(音が鳴り始める)
			a.ch3.enabled = a.ch3.dacEnable
			if a.ch3.length == 0 {
				a.ch3.length = 256
			}
			a.ch3.window = 0
			if a.model == APU_GBA {
				if a.ch3.mode == 1 {
					a.ch3.usedBank = 0
				} else {
					a.ch3.usedBank = a.ch3.bank
				}
			}
		}

	case 0xFF20:
		a.ch4.length = 64 - int(val&0b11_1111)
	case 0xFF21:
		a.ch4.envelope.initialVolume = int(val>>4) & 0b1111
		a.ch4.envelope.direction = util.Bit(val, 3)
		a.ch4.envelope.speed = int(val & 0b111)
		if !a.ch4.dacEnable() {
			a.ch4.enabled = false
		}
	case 0xFF22:
		a.ch4.octave = int(val >> 4) // ノイズ周波数2(オクターブ指定)
		a.ch4.divisor = int(val & 0b111)
		a.ch4.period = a.ch4.calcFreqency()
		a.ch4.width = 15
		if util.Bit(val, 3) {
			a.ch4.width = 7
		}
	case 0xFF23:
		a.ch4.stop = util.Bit(val, 6)
		if util.Bit(val, 7) { // キーオン(音が鳴り始める)
			a.ch4.tryRestart()
		}

	case 0xFF24:
		a.volume[0] = int(val>>4) & 0b111
		a.volume[1] = int(val & 0b111)

	case 0xFF25:
		a.ch1.ignored = !util.Bit(val, 0)
		a.ch2.ignored = !util.Bit(val, 1)
		a.ch3.ignored = !util.Bit(val, 2)
		a.ch4.ignored = !util.Bit(val, 3)

	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		if a.model == APU_GBA {
			if a.ch3.bank == 0 {
				a.ch3.samples[16+addr-0xFF30] = val
			} else {
				a.ch3.samples[addr-0xFF30] = val
			}
		} else {
			a.ch3.samples[addr-0xFF30] = val
		}
	}
}
