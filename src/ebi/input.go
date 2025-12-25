package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var KeyBindings = map[string]ebiten.Key{
	"A":      ebiten.KeyX,
	"B":      ebiten.KeyZ,
	"SELECT": ebiten.KeyBackspace,
	"START":  ebiten.KeyEnter,
	"UP":     ebiten.KeyArrowUp,
	"DOWN":   ebiten.KeyArrowDown,
	"LEFT":   ebiten.KeyArrowLeft,
	"RIGHT":  ebiten.KeyArrowRight,
}

var GamepadBindings = map[string]ebiten.StandardGamepadButton{
	"A":      ebiten.StandardGamepadButtonRightRight,
	"B":      ebiten.StandardGamepadButtonRightBottom,
	"SELECT": ebiten.StandardGamepadButtonCenterLeft,
	"START":  ebiten.StandardGamepadButtonCenterRight,
	"UP":     ebiten.StandardGamepadButtonLeftTop,
	"DOWN":   ebiten.StandardGamepadButtonLeftBottom,
	"LEFT":   ebiten.StandardGamepadButtonLeftLeft,
	"RIGHT":  ebiten.StandardGamepadButtonLeftRight,
}

var Inputs = map[string]bool{
	"A":      false,
	"B":      false,
	"START":  false,
	"SELECT": false,
	"UP":     false,
	"DOWN":   false,
	"LEFT":   false,
	"RIGHT":  false,
}

func Input() {
	for key := range Inputs {
		Inputs[key] = false
	}

	// State save/load (same with RetroArch)
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		App.Emu.SaveState()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		App.Emu.LoadState()
	}
	pollKeyboard()
	pollGamepad()
}

func pollKeyboard() {
	for input, key := range KeyBindings {
		if _, ok := Inputs[input]; ok {
			if ebiten.IsKeyPressed(key) {
				Inputs[input] = true
			}
		}
	}
}

func pollGamepad() {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {

		for input, b := range GamepadBindings {
			switch {
			case ebiten.IsStandardGamepadButtonPressed(id, b):
				if _, ok := Inputs[input]; ok {
					Inputs[input] = true
				}
			}
		}

		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) {
		case 1:
			Inputs["RIGHT"] = true
		case -1:
			Inputs["LEFT"] = true
		}
		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) {
		case 1:
			Inputs["DOWN"] = true
		case -1:
			Inputs["UP"] = true
		}
	}
}
