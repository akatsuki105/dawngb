package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var keyMap = map[ebiten.Key]string{
	ebiten.KeyX:          "A",
	ebiten.KeyZ:          "B",
	ebiten.KeyBackspace:  "SELECT",
	ebiten.KeyEnter:      "START",
	ebiten.KeyArrowUp:    "UP",
	ebiten.KeyArrowDown:  "DOWN",
	ebiten.KeyArrowLeft:  "LEFT",
	ebiten.KeyArrowRight: "RIGHT",
}

var standardButtonToString = map[ebiten.StandardGamepadButton]string{
	ebiten.StandardGamepadButtonRightTop:    "B",
	ebiten.StandardGamepadButtonRightRight:  "A",
	ebiten.StandardGamepadButtonCenterLeft:  "SELECT",
	ebiten.StandardGamepadButtonCenterRight: "START",
	ebiten.StandardGamepadButtonLeftBottom:  "DOWN",
	ebiten.StandardGamepadButtonLeftRight:   "RIGHT",
	ebiten.StandardGamepadButtonLeftLeft:    "LEFT",
	ebiten.StandardGamepadButtonLeftTop:     "UP",
}

var inputMap = map[string]bool{
	"A":      false,
	"B":      false,
	"START":  false,
	"SELECT": false,
	"UP":     false,
	"DOWN":   false,
	"LEFT":   false,
	"RIGHT":  false,
}

var inputMapWeb = map[string]bool{
	"A":      false,
	"B":      false,
	"START":  false,
	"SELECT": false,
	"UP":     false,
	"DOWN":   false,
	"LEFT":   false,
	"RIGHT":  false,
}

func (e *Emu) pollInput() {
	for key := range inputMap {
		inputMap[key] = inputMapWeb[key]
	}

	e.pollKeyInput()
	e.pollGamepadInput()

	for key, input := range inputMap {
		e.c.SetKeyInput(key, input)
	}

	if e.c != nil {
		if ebiten.IsKeyPressed(ebiten.KeyDigit1) && !saved {
			fmt.Println("State save")
			ok := e.c.Serialize(&stateSaveData)
			if !ok {
				fmt.Println("State save failed")
			}
			// fmt.Println(stateSaveData.Len())
			saved = true
		} else if ebiten.IsKeyPressed(ebiten.KeyDigit0) && saved {
			fmt.Println("State load")
			ok := e.c.Deserialize(&stateSaveData)
			if !ok {
				fmt.Println("State load failed")
			}
			// fmt.Println(stateSaveData.Len())
			saved = false
		}
	}
}

func (e *Emu) pollKeyInput() {
	for key, input := range keyMap {
		if _, ok := inputMap[input]; ok {
			if ebiten.IsKeyPressed(key) {
				inputMap[input] = true
			}
		}
	}
}

func (e *Emu) pollGamepadInput() {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {

		for b, input := range standardButtonToString {
			switch {
			case ebiten.IsStandardGamepadButtonPressed(id, b):
				if _, ok := inputMap[input]; ok {
					inputMap[input] = true
				}
			}
		}

		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) {
		case 1:
			inputMap["RIGHT"] = true
		case -1:
			inputMap["LEFT"] = true
		}
		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) {
		case 1:
			inputMap["DOWN"] = true
		case -1:
			inputMap["UP"] = true
		}
	}
}
