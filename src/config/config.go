package config

type GB struct {
	Model uint8 // 0: DMG, 1: SGB, 2: CGB
	Intro bool
}

type Audio struct {
	Enable bool
	Volume float64
}

type Logger struct {
	Enable bool
	Level  string
}

type Config struct {
	ShowFPS bool
	Audio   Audio
	Logger  Logger
	GB      GB
}

var DefaultConfig = Config{
	ShowFPS: false,
	Audio: Audio{
		Enable: true,
		Volume: 0.5,
	},
	Logger: Logger{
		Enable: true,
		Level:  "info",
	},
	GB: GB{
		Model: 2,
		Intro: true,
	},
}
