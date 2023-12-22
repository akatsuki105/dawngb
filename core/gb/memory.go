package gb

type Memory struct {
	gb *GB
}

func newMemory(gb *GB) *Memory {
	return &Memory{
		gb: gb,
	}
}

func (m *Memory) Read(addr uint16) byte {
	return 0
}

func (m *Memory) Write(addr uint16, val byte) {}
