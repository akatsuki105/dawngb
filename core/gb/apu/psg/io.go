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

func (a *PSG) Read(addr uint16, peek bool) uint8 {
	switch addr {
	case NR52:
		val := uint8(0)
		val = setBit(val, 7, a.Enabled)
		val = setBit(val, 0, a.CH1.enabled)
		val = setBit(val, 1, a.CH2.enabled)
		val = setBit(val, 2, a.CH3.enabled)
		val = setBit(val, 3, a.CH4.enabled)
		return val
	}
	if addr >= 0xFF30 && addr < 0xFF40 {
		if peek {
			return a.CH3.Peek(addr)
		}
		return a.CH3.read(addr)
	}

	return a.ioreg[addr-0xFF10]
}

func (a *PSG) Write(addr uint16, val uint8) {
	if addr == NR52 {
		prev := a.Enabled
		a.Enabled = getBit(val, 7)
		if prev && !a.Enabled { // APUがオンからオフになったとき
			a.rightVolume, a.leftVolume = 0, 0
			a.rightEnables[0], a.rightEnables[1], a.rightEnables[2], a.rightEnables[3] = false, false, false, false
			a.leftEnables[0], a.leftEnables[1], a.leftEnables[2], a.leftEnables[3] = false, false, false, false
			a.CH1.TurnOff()
			a.CH2.TurnOff()
			a.CH3.TurnOff()
			a.CH4.TurnOff()
		}
		return
	}

	if !a.Enabled {
		return
	}

	a.ioreg[addr-0xFF10] = val
	switch addr {
	case NR10:
		negate := getBit(val, 3)
		if a.CH1.sweep.enabled && (a.CH1.sweep.negate && !negate) {
			a.CH1.enabled = false // 下降スイープ中に上昇スイープを設定すると消音される
		}
		a.CH1.sweep.shift = val & 0b111
		a.CH1.sweep.negate = negate
		a.CH1.sweep.interval = (val >> 4) & 0b111
		return
	case NR11:
		a.CH1.length = 64 - (val & 0b11_1111)
		a.CH1.duty = (val >> 6)
		return
	case NR12:
		a.CH1.envelope.speed = (val & 0b111)
		a.CH1.envelope.direction = getBit(val, 3)
		a.CH1.envelope.initialVolume = (val >> 4) & 0b1111
		if !a.CH1.dacEnable() {
			a.CH1.enabled = false
		}
		return
	case NR13:
		a.CH1.period = (a.CH1.period & 0xFF00) | uint16(val)
		return
	case NR14:
		a.CH1.period = (a.CH1.period & 0x00FF) | (uint16(val&0b111) << 8)
		a.CH1.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.CH1.reload()
		}
		return

	case NR21:
		a.CH2.length = 64 - (val & 0b11_1111)
		a.CH2.duty = (val >> 6)
		return
	case NR22:
		a.CH2.envelope.speed = (val & 0b111)
		a.CH2.envelope.direction = getBit(val, 3)
		a.CH2.envelope.initialVolume = (val >> 4) & 0b1111
		if !a.CH2.dacEnable() {
			a.CH2.enabled = false
		}
		return
	case NR23:
		a.CH2.period = (a.CH2.period & 0xFF00) | uint16(val)
		return
	case NR24:
		a.CH2.period = (a.CH2.period & 0x00FF) | (uint16(val&0b111) << 8)
		a.CH2.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.CH2.reload()
		}
		return

	case NR30:
		a.CH3.dacEnable = getBit(val, 7)
		if !a.CH3.dacEnable {
			a.CH3.enabled = false
		}
		if a.model == MODEL_GBA {
			a.CH3.mode = (val >> 5) & 0b1
			a.CH3.Bank = (val >> 6) & 0b1
		}
		return
	case NR31:
		a.CH3.length = 256 - uint16(val)
		return
	case NR32:
		a.CH3.volume = (val >> 5) & 0b11
		return
	case NR33:
		a.CH3.period &= 0b111_0000_0000
		a.CH3.period |= uint16(val)
		a.CH3.freqCounter = a.CH3.windowStepCycle()
		return
	case NR34:
		a.CH3.stop = getBit(val, 6)
		a.CH3.period &= 0b000_1111_1111
		a.CH3.period |= uint16(val&0b111) << 8
		a.CH3.freqCounter = a.CH3.windowStepCycle()
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.CH3.reload()
			if a.model == MODEL_GBA {
				if a.CH3.mode == 1 {
					a.CH3.curBank = 0
				} else {
					a.CH3.curBank = a.CH3.Bank
				}
			}
		}
		return

	case NR41:
		a.CH4.length = 64 - (val & 0b11_1111)
		return
	case NR42:
		a.CH4.envelope.initialVolume = (val >> 4) & 0b1111
		a.CH4.envelope.direction = getBit(val, 3)
		a.CH4.envelope.speed = (val & 0b111)
		if !a.CH4.dacEnable() {
			a.CH4.enabled = false
		}
		return
	case NR43:
		a.CH4.divisor = (val & 0b111)
		a.CH4.narrow = getBit(val, 3)
		a.CH4.octave = (val >> 4) // ノイズ周波数2(オクターブ指定)
		a.CH4.period = a.CH4.calcFreqency()
		return
	case NR44:
		a.CH4.stop = getBit(val, 6)
		if getBit(val, 7) { // キーオン(音が鳴り始める)
			a.CH4.reload()
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
		a.CH3.write(addr, val)
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
