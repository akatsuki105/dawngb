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

	points := NewPackedVector2Array()
	points.Resize(160)

	colors := NewPackedColorArray()
	colors.Resize(160)

	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			r, g, b := float32(data[y*160+x].R)/255.0, float32(data[y*160+x].G)/255.0, float32(data[y*160+x].B)/255.0
			points.SetIndexed(int64(x), NewVector2WithFloat32Float32(float32(x)*SCALE, float32(y)*SCALE))
			colors.SetIndexed(int64(x), NewColorWithFloat32Float32Float32(r, g, b))
			s.DrawRect(NewRect2WithFloat32Float32Float32Float32(float32(x)*SCALE, float32(y)*SCALE, SCALE, SCALE), NewColorWithFloat32Float32Float32(r, g, b), false, -1)
		}
		// s.getFrameBuffer().DrawMultilineColors(points, colors, -1)
		// s.DrawPolygon(points, colors, NewPackedVector2Array(), nil)
	}
	// s.DrawTextureRect(s.getFrameBuffer().GetTexture(), NewRect2WithFloat32Float32Float32Float32(0, 0, 160*SCALE, 144*SCALE), false, NewColorWithFloat32Float32Float32(1, 1, 1), false)
}

func (s *Screen) V_Process(delta float32) {
	gb.RunFrame()

	if counter%3 == 0 {
		s.QueueRedraw()
	}
	counter++
}

func (s *Screen) getFrameBuffer() TextureRect {
	gds := NewStringWithLatin1Chars("FrameBuffer")
	defer gds.Destroy()
	path := NewNodePathWithString(gds)
	defer path.Destroy()
	return ObjectCastTo(s.GetNode(path), "TextureRect").(TextureRect)
}
