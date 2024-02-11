package scene

import (
	_ "embed"

	"github.com/akatsuki105/dawngb/core"
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/gdclassimpl"
)

const (
	SCALE = 1.
)

var counter = 0

//go:embed hello.gb
var rom []byte

var gb core.Core

// 初期化用の関数
func RegisterClassScreen() {
	ClassDBRegisterClass[*Screen](&Screen{}, []GDExtensionPropertyInfo{}, nil, func(t GDClass) {
		// V_Readyと_ready、V_Processと_processを紐付ける
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Draw", "_draw", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Process", "_process", []string{"delta"}, nil)
	})
}

type Screen struct {
	gdclassimpl.ControlImpl
}

func (s *Screen) GetClassName() string {
	return "Screen"
}

func (s *Screen) GetParentClassName() string {
	return "Control"
}

func (s *Screen) V_Ready() {
	gb = core.New("GB", nil)
	if err := gb.LoadROM(rom); err != nil {
		panic(err)
	}
}

func (s *Screen) V_Draw() {
	data := gb.Screen()
	if len(data) != 0 {
		for i := 0; i < (160 * 144); i++ {
			x := i % 160
			y := i / 160
			r, g, b := float32(data[i].R)/255.0, float32(data[i].G)/255.0, float32(data[i].B)/255.0
			s.DrawRect(NewRect2WithFloat32Float32Float32Float32(float32(x)*SCALE, float32(y)*SCALE, SCALE, SCALE), NewColorWithFloat32Float32Float32(r, g, b), true, -1.0)
		}
	}
}

func (s *Screen) V_Process(delta float32) {
	gb.RunFrame()
	if counter%3 == 0 {
		s.QueueRedraw()
	}
	counter++
}
