package scene

import (
	"github.com/godot-go/godot-go/pkg/builtin"
	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// 初期化用の関数
func RegisterClassScreen() {
	core.ClassDBRegisterClass(&Screen{}, []ffi.GDExtensionPropertyInfo{}, nil, func(t builtin.GDClass) {
		// V_Readyと_ready、V_Processと_processを紐付ける
		core.ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		core.ClassDBBindMethodVirtual(t, "V_Process", "_process", []string{"delta"}, nil)
	})
}

type Screen struct {
	gdclassimpl.ControlImpl
	input builtin.Input
}

func (*Screen) GetClassName() string {
	return "Screen"
}

func (*Screen) GetParentClassName() string {
	return "Control"
}

func (s *Screen) V_Ready() {}

func (s *Screen) V_Process(delta float32) {}
