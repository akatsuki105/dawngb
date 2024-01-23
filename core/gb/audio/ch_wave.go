package audio

type wave struct {
	enabled bool
	ignored bool // Ignore sample output

	dacEnable bool // NR30's bit7
	stop      bool // .length が 0 になったときに 音を止めるかどうか(NR34's bit6)
	length    int  // 音の残り再生時間
	volume    int  // NR32's bit6-5 (0: 0%, 1: 100%, 2: 50%, 3: 25%)

	period      int // GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter int

	samples [16]uint8 // 4bit sample
	window  int       // 0 ~ 31
}

func newWaveChannel() *wave {
	return &wave{
		ignored: true,
	}
}

func (ch *wave) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length <= 0 {
			ch.enabled = false
		}
	}
}

func (ch *wave) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
	} else {
		ch.freqCounter = ch.windowStepCycle()
		ch.window = (ch.window + 1) & 0x1F
	}
}

func (ch *wave) getOutput() int {
	if !ch.ignored {
		if ch.enabled && ch.dacEnable {
			isHi := ch.window&1 == 0 // 上位4bit -> 下位4bit -> 上位4bit -> 下位4bit -> ...
			sample := uint8(0)
			if isHi {
				sample = ch.samples[ch.window>>1] >> 4
			} else {
				sample = ch.samples[ch.window>>1] & 0xF
			}
			return int(sample >> ch.volume)
		}
	}
	return 0
}

func (ch *wave) windowStepCycle() int {
	return 2 * (2048 - ch.period)
}
