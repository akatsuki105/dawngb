package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/akatsuki105/dawngb/src/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// ExitCode represents program's status code
type ExitCode int

// exit code
const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
)

type AppState struct {
	Name    string // アプリケーション名
	Counter int    // フレームカウンタ
	Config  config.Config
	Emu     *Emu
	Audio   *AudioManager
	Logger  *slog.Logger
}

var App = AppState{
	Name:   "DawnGB",
	Config: config.DefaultConfig,
}

func main() {
	os.Exit(int(Run()))
}

func Run() ExitCode {
	flag.Parse()

	App.initLogger()

	App.Audio = NewAudioManager()
	defer App.Audio.Close()

	App.Emu = createEmu(App.Config.GB.Model)

	if flag.NArg() > 0 {
		err := App.Emu.LoadROMFromPath(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}
	}

	ebiten.SetTPS(60)
	ebiten.SetWindowTitle(fmt.Sprintf("%s - 60.0 FPS", App.Name))
	ebiten.SetScreenClearedEveryFrame(false)

	{
		f := ebiten.Monitor().DeviceScaleFactor()
		w, h := float64(160)*f, float64(144)*f
		ebiten.SetWindowSize(int(w), int(h))
	}

	if err := ebiten.RunGame(&App); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func (app *AppState) Update() error {
	if app.Config.Audio.Enable {
		app.Audio.Update()
	}

	Input()
	app.Emu.Update()
	err := app.HandleDropFile()
	if err != nil {
		slog.Error("Failed to handle dropped file", "error", err)
	}

	app.Counter++
	if app.Counter&0xFF == 0 {
		fps := ebiten.ActualTPS()
		title := fmt.Sprintf("%s - %.1f FPS", app.Name, fps)
		ebiten.SetWindowTitle(title)
	}
	return nil
}

func (app *AppState) Draw(screen *ebiten.Image) {
	App.Emu.Draw(screen)
}

// 引数にウィンドウサイズをとり、画面の解像度を返す
func (app *AppState) Layout(_, _ int) (screenWidth, screenHeight int) { return 160, 144 }

func (app *AppState) initLogger() {
	cfg := &app.Config.Logger
	if cfg.Enable {
		level := slog.LevelInfo
		switch strings.ToLower(cfg.Level) {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
		app.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	} else {
		app.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	slog.SetDefault(app.Logger)
}

func (app *AppState) HandleDropFile() error {
	file := ebiten.DroppedFiles()
	if file != nil {
		entries, err := fs.ReadDir(file, ".")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				data, err := fs.ReadFile(file, name)
				if err != nil {
					return err
				}

				switch filepath.Ext(name) {
				case ".gb", ".gbc": // ROM
					err := App.Emu.LoadROM(data)
					if err != nil {
						return err
					}

				case ".sav", ".srm": // Save Data
					err := App.Emu.LoadSave(data)
					if err != nil {
						return err
					}

				case ".bin": // BIOS
					err := App.Emu.LoadBIOS(data)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
