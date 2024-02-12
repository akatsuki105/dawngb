package main

import "C"
import (
	"unsafe"

	"github.com/akatsuki105/dawngb/src/godot/scene"
	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
)

//export GameBoyInit
func GameBoyInit(p_get_proc_address unsafe.Pointer, p_library unsafe.Pointer, r_initialization unsafe.Pointer) bool {

	initObj := core.NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(p_get_proc_address),
		(ffi.GDExtensionClassLibraryPtr)(p_library),
		(*ffi.GDExtensionInitialization)(unsafe.Pointer(r_initialization)),
	)

	initObj.RegisterSceneInitializer(func() {
		scene.RegisterClassScreen()
	})

	initObj.RegisterSceneTerminator(func() {
	})

	return initObj.Init()
}

func main() {
}
