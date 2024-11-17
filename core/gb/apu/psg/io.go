package psg

// NOTE: "ゼルダの伝説 ふしぎの木の実" はAPUのNR52を正しく実装しないとタイトル画面から進めない

const (
	// CH1
	NR10 = 0xFF10
	NR11 = 0xFF11
	NR12 = 0xFF12
	NR13 = 0xFF13
	NR14 = 0xFF14

	// CH2
	NR20 = 0xFF15 // CH2にはスイープ機能がないのでレジスタは存在しない
	NR21 = 0xFF16
	NR22 = 0xFF17
	NR23 = 0xFF18
	NR24 = 0xFF19

	// CH3
	NR30 = 0xFF1A
	NR31 = 0xFF1B
	NR32 = 0xFF1C
	NR33 = 0xFF1D
	NR34 = 0xFF1E

	// CH4
	NR40 = 0xFF1F
	NR41 = 0xFF20
	NR42 = 0xFF21
	NR43 = 0xFF22
	NR44 = 0xFF23

	// Control
	NR50 = 0xFF24
	NR51 = 0xFF25
	NR52 = 0xFF26
)

func (a *PSG) Read(addr uint16) uint8 {
	switch addr {
	case NR52:
		val := uint8(0)
		val = setBit(val, 7, a.enabled)
		val = setBit(val, 0, a.ch1.enabled)
		val = setBit(val, 1, a.ch2.enabled)
		val = setBit(val, 2, a.ch3.enabled)
		val = setBit(val, 3, a.ch4.enabled)
		return val
	}

	if addr >= 0xFF30 && addr < 0xFF40 {
		return a.ch3.read(addr)
	}

	return a.ioreg[addr-0xFF10]
}

func (a *PSG) Write(addr uint16, val uint8) {
	if addr == NR52 {
		a.enabled = getBit(val, 7)
		return
	}

	if !a.enabled {
		return
	}

	a.ioreg[addr-0xFF10] = val
	switch addr {
	case NR10:
		negate := getBit(val, 3)
		if a.ch1.sweep.enabled && (a.ch1.sweep.negate && !negate) {
			a.ch1.enabled = false // 下降スイープ中に上昇スイープを設定すると消音される
		}
		a.ch1.sweep.shift = val & 0b111
		a.ch1.sweep.negate = negate
		a.ch1.sweep.interval = (val >> 4) & 0b111
		return
	case NR11:
		a.ch1.length = 64 - (val & 0b11_1111)
		a.ch1.duty = (val >> 6)
		return
	case NR12:
		a.ch1.envelope.speed = (val & 0b111)
		a.ch1.envelope.direction = getBit(val, 3)
		a.ch1.envelope.initialVolume = (val >> 4) & 0b1111
		if !a.ch1.dacEnable() {
			a.ch1.enabled = false
		}
		return
	case NR13:
		a.ch1.period = (a.ch1.period & 0xFF00) | uint16(val)
		return
	case NR14:
		a.ch1.period = (a.ch1.period & 0x00FF) | (uint16(val&0b111) << 8)
		a.ch1.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.ch1.reload()
		}
		return

	case NR21:
		a.ch2.length = 64 - (val & 0b11_1111)
		a.ch2.duty = (val >> 6)
		return
	case NR22:
		a.ch2.envelope.speed = (val & 0b111)
		a.ch2.envelope.direction = getBit(val, 3)
		a.ch2.envelope.initialVolume = (val >> 4) & 0b1111
		if !a.ch2.dacEnable() {
			a.ch2.enabled = false
		}
		return
	case NR23:
		a.ch2.period = (a.ch2.period & 0xFF00) | uint16(val)
		return
	case NR24:
		a.ch2.period = (a.ch2.period & 0x00FF) | (uint16(val&0b111) << 8)
		a.ch2.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.ch2.reload()
		}
		return

	case NR30:
		a.ch3.dacEnable = getBit(val, 7)
		if !a.ch3.dacEnable {
			a.ch3.enabled = false
		}
		if a.model == MODEL_GBA {
			a.ch3.mode = (val >> 5) & 0b1
			a.ch3.bank = (val >> 6) & 0b1
		}
		return
	case NR31:
		a.ch3.length = 256 - uint16(val)
		return
	case NR32:
		a.ch3.volume = (val >> 5) & 0b11
		return
	case NR33:
		a.ch3.period &= 0b111_0000_0000
		a.ch3.period |= uint16(val)
		a.ch3.freqCounter = a.ch3.windowStepCycle()
		return
	case NR34:
		a.ch3.stop = getBit(val, 6)
		a.ch3.period &= 0b000_1111_1111
		a.ch3.period |= uint16(val&0b111) << 8
		a.ch3.freqCounter = a.ch3.windowStepCycle()
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.ch3.reload()
			if a.model == MODEL_GBA {
				if a.ch3.mode == 1 {
					a.ch3.curBank = 0
				} else {
					a.ch3.curBank = a.ch3.bank
				}
			}
		}
		return

	case NR41:
		a.ch4.length = 64 - (val & 0b11_1111)
		return
	case NR42:
		a.ch4.envelope.initialVolume = (val >> 4) & 0b1111
		a.ch4.envelope.direction = getBit(val, 3)
		a.ch4.envelope.speed = (val & 0b111)
		if !a.ch4.dacEnable() {
			a.ch4.enabled = false
		}
		return
	case NR43:
		a.ch4.divisor = (val & 0b111)
		a.ch4.narrow = getBit(val, 3)
		a.ch4.octave = (val >> 4) // ノイズ周波数2(オクターブ指定)
		a.ch4.period = a.ch4.calcFreqency()
		return
	case NR44:
		a.ch4.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.ch4.reload()
		}
		return

	case NR50:
		a.rightVolume = (val >> 0) & 0b111
		a.leftVolume = (val >> 4) & 0b111
		return
	case NR51:
		a.rightEnables[0] = getBit(val, 0)
		a.rightEnables[1] = getBit(val, 1)
		a.rightEnables[2] = getBit(val, 2)
		a.rightEnables[3] = getBit(val, 3)
		a.leftEnables[0] = getBit(val, 4)
		a.leftEnables[1] = getBit(val, 5)
		a.leftEnables[2] = getBit(val, 6)
		a.leftEnables[3] = getBit(val, 7)
		return
	}

	if addr >= 0xFF30 && addr < 0xFF40 {
		a.ch3.write(addr, val)
		return
	}
}

func getBit(val uint8, bit int) bool {
	return val&(1<<bit) != 0
}

func setBit(val uint8, bit int, b bool) uint8 {
	if b {
		return val | (1 << bit)
	}
	return val & ^(1 << bit)
}
