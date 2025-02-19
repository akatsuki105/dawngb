package psg

const waveBank = 2

// .volume が index
var volumeShift = [4]uint8{4, 0, 1, 2} // 波形は最大15なので4左シフトすれば0%

// 波形メモリ音源
type wave struct {
	model   uint8
	enabled bool // NR52.2

	dacEnable bool   // NR30.7
	length    uint16 // NR31; 音の残り再生時間
	volume    uint8  // NR32.5-6; 0: 0%, 1: 100%, 2: 50%, 3: 25%
	stop      bool   // NR34.6; .length が 0 になったときに 音を止めるかどうか

	period      uint16 // NR33.0-7, NR34.0-2; GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter uint16

	RAM    [16 * waveBank]uint8 // 4bitサンプル*32 で16バイト ; GBAの場合はバンクが2つある
	sample uint8                // 0..15
	window uint8                // 0..31

	output uint8 // 0..15

	// For GBA
	mode    uint8 // NR30.5; 0: 16バイト(32サンプル)を演奏に使い、裏のバンクでは読み書きを行う、 1: 32バイト(64サンプル)を全部演奏に使う
	Bank    uint8 // NR30.6
	curBank uint8 // 現在演奏中のバンク、modeが1の場合は、 .bank の値と必ずしも一致しないので
}

func newWaveChannel(model uint8) *wave {
	return &wave{
		model: model,
	}
}

func (ch *wave) Reset() {
	ch.TurnOff()
	ch.period, ch.freqCounter = 0, 0
	ch.sample, ch.window = 0, 0
	clear(ch.RAM[:])
	ch.output = 0
	ch.mode, ch.Bank, ch.curBank = 0, 0, 0
}

func (ch *wave) TurnOff() {
	ch.dacEnable = false
	ch.length, ch.volume, ch.stop, ch.period = 0, 0, false, 0
	ch.enabled = false
}

func (ch *wave) reload() {
	ch.enabled = ch.dacEnable
	ch.freqCounter = ch.windowStepCycle() + 2
	ch.window = 0
	ch.output = 0
	if ch.length == 0 {
		ch.length = 256
	}
}

func (ch *wave) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *wave) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
		if ch.freqCounter == 0 {
			ch.freqCounter = ch.windowStepCycle()
			ch.update()
		}
	}
}

// GetOutput gets 4bit sample (0..15)
func (ch *wave) GetOutput() uint8 {
	if ch.enabled {
		shift := volumeShift[ch.volume]
		return ch.output >> shift
	}
	return 0
}

func (ch *wave) update() {
	ch.window = (ch.window + 1) & 0x1F // 読み出す前にインクリメント(CH3のreload後に最初に読み出すのはsamples[0]の下位ニブル)

	upper := (ch.window & 0x1) == 0
	if upper {
		ch.output = ch.RAM[ch.window>>1] >> 4
	} else {
		ch.output = ch.RAM[ch.window>>1] & 0xF
	}

	if ch.window == 0 {
		ch.curBank ^= ch.mode
	}
}

func (ch *wave) read(addr uint16) uint8 {
	if !ch.enabled {
		bank := uint16(0)
		if ch.model == MODEL_GBA {
			if ch.Bank == 0 {
				bank = 16
			}
		}
		return ch.RAM[bank|(addr&0xF)]
	}
	return 0xFF // AGB
}

func (ch *wave) write(addr uint16, val uint8) {
	if !ch.enabled {
		bank := uint16(0)
		if ch.model == MODEL_GBA {
			if ch.Bank == 0 {
				bank = 16
			}
		}
		ch.RAM[bank|(addr&0xF)] = val
	}
}

func (ch *wave) Peek(addr uint16) uint8 {
	bank := uint16(0)
	if ch.model == MODEL_GBA {
		bank = uint16(ch.curBank) * 16
	}
	return ch.RAM[bank|(addr&0xF)]
}

func (ch *wave) windowStepCycle() uint16 {
	return 2 * (2048 - ch.period)
}

type WaveSnapshot struct {
	Header              uint64
	Enabled             bool
	DAC                 bool
	Length              uint16
	Volume              uint8
	Stop                bool
	Period, FreqCounter uint16
	RAM                 [16 * waveBank]uint8
	Sample, Window      uint8
	Output              uint8
	Mode, Bank, CurBank uint8
	Reserved            [16]uint8
}

func (ch *wave) CreateSnapshot() WaveSnapshot {
	snap := WaveSnapshot{
		Enabled:     ch.enabled,
		DAC:         ch.dacEnable,
		Length:      ch.length,
		Volume:      ch.volume,
		Stop:        ch.stop,
		Period:      ch.period,
		FreqCounter: ch.freqCounter,
		Sample:      ch.sample,
		Window:      ch.window,
		Output:      ch.output,
		Mode:        ch.mode,
		Bank:        ch.Bank,
		CurBank:     ch.curBank,
	}
	copy(snap.RAM[:], ch.RAM[:])
	return snap
}

func (ch *wave) RestoreSnapshot(snap WaveSnapshot) bool {
	ch.enabled = snap.Enabled
	ch.dacEnable = snap.DAC
	ch.length, ch.volume, ch.stop, ch.period = snap.Length, snap.Volume, snap.Stop, snap.Period
	ch.freqCounter = snap.FreqCounter
	ch.sample, ch.window = snap.Sample, snap.Window
	ch.output = snap.Output
	ch.mode, ch.Bank, ch.curBank = snap.Mode, snap.Bank, snap.CurBank
	copy(ch.RAM[:], snap.RAM[:])
	return true
}
