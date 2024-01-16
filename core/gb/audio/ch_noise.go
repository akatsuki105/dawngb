package audio

type noise struct {
	ignored bool // Ignore sample output
}

func newNoiseChannel() *noise {
	return &noise{
		ignored: true,
	}
}

func (ch *noise) getOutput() int {
	if !ch.ignored {
		return 0
	}
	return 0
}
