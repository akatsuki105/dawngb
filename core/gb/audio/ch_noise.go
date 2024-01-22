package audio

type noise struct {
	enabled bool
	ignored bool // Ignore sample output

	length int // 音の残り再生時間
	stop   bool

	envelope *envelope

	lfsr uint16 // シフトレジスタ(Noiseの疑似乱数)

	// この2つでノイズの周波数(乱数の生成頻度)を決める
	octave  int // ノイズ周波数2(オクターブ指定)
	divisor int // ノイズ周波数1(カウント指定)
	period  int

	width int

	output int
}

func newNoiseChannel() *noise {
	return &noise{
		ignored:  true,
		envelope: newEnvelope(),
		lfsr:     1,
		width:    15,
	}
}

func (ch *noise) clock64Hz() {
	if ch.enabled {
		ch.enabled = ch.envelope.update()
	}
}

func (ch *noise) clock256Hz() {
	if ch.enabled {
		if ch.stop && ch.length > 0 {
			ch.length--
			if ch.length <= 0 {
				ch.enabled = false
			}
		}
	}
}

func (ch *noise) clockTimer() {
	// ch.enabledに関わらず、乱数は生成される
	result := 0
	if ch.enabled {
		ch.period--
		if ch.period <= 0 {
			ch.period = ch.calcFreqency()
			if ch.octave < 14 {
				mask := ((ch.lfsr ^ (ch.lfsr >> 1)) & 1)
				ch.lfsr = ((ch.lfsr >> 1) ^ (mask << (ch.width - 1))) & 0x7FFF
			}
		}

		if (ch.lfsr & 1) == 0 {
			result = ch.envelope.volume
		}
	}

	ch.output = result
}

func (ch *noise) calcFreqency() int {
	freq := 1
	if ch.divisor != 0 {
		freq = 2 * ch.divisor
	}
	freq <<= ch.octave
	return freq * 8
}

func (ch *noise) getOutput() int {
	if !ch.ignored {
		if ch.enabled {
			return ch.output
		}
	}
	return 0
}
