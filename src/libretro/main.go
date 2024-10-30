package main

// #include "./libretro.h"
import "C"
import (
	"bytes"
	"os"
	"unsafe"

	"github.com/akatsuki105/dawngb/core"
)

const (
	width  = 160
	height = 144
)

const (
	RETRO_REGION_NTSC = 0
	RETRO_REGION_PAL  = 1
)

const (
	RETRO_ENVIRONMENT_SET_PIXEL_FORMAT   = 10
	RETRO_ENVIRONMENT_GET_INPUT_BITMASKS = 51 | RETRO_ENVIRONMENT_EXPERIMENTAL
	RETRO_ENVIRONMENT_EXPERIMENTAL       = 0x10000
)

var (
	useBitmasks bool
	keymap      = []uint{
		C.RETRO_DEVICE_ID_JOYPAD_A,
		C.RETRO_DEVICE_ID_JOYPAD_B,
		C.RETRO_DEVICE_ID_JOYPAD_SELECT,
		C.RETRO_DEVICE_ID_JOYPAD_START,
		C.RETRO_DEVICE_ID_JOYPAD_RIGHT,
		C.RETRO_DEVICE_ID_JOYPAD_LEFT,
		C.RETRO_DEVICE_ID_JOYPAD_UP,
		C.RETRO_DEVICE_ID_JOYPAD_DOWN,
	}
	keymapNames = map[uint]string{
		C.RETRO_DEVICE_ID_JOYPAD_A:      "A",
		C.RETRO_DEVICE_ID_JOYPAD_B:      "B",
		C.RETRO_DEVICE_ID_JOYPAD_SELECT: "SELECT",
		C.RETRO_DEVICE_ID_JOYPAD_START:  "START",
		C.RETRO_DEVICE_ID_JOYPAD_RIGHT:  "RIGHT",
		C.RETRO_DEVICE_ID_JOYPAD_LEFT:   "LEFT",
		C.RETRO_DEVICE_ID_JOYPAD_UP:     "UP",
		C.RETRO_DEVICE_ID_JOYPAD_DOWN:   "DOWN",
	}
)

type emulator struct {
	c            core.Core
	samples      []byte
	sampleBuffer *bytes.Buffer
}

var e *emulator

//export retro_set_environment
func retro_set_environment(cb C.retro_environment_t) {
	C._retro_set_environment(cb)
}

//export retro_set_video_refresh
func retro_set_video_refresh(cb C.retro_video_refresh_t) {
	C._retro_set_video_refresh(cb)
}

//export retro_set_audio_sample
func retro_set_audio_sample(cb C.retro_audio_sample_t) {
	C._retro_set_audio_sample(cb)
}

//export retro_set_audio_sample_batch
func retro_set_audio_sample_batch(cb C.retro_audio_sample_batch_t) {
	C._retro_set_audio_sample_batch(cb)
}

//export retro_set_input_poll
func retro_set_input_poll(cb C.retro_input_poll_t) {
	C._retro_set_input_poll(cb)
}

//export retro_set_input_state
func retro_set_input_state(cb C.retro_input_state_t) {
	C._retro_set_input_state(cb)
}

//export retro_init
func retro_init() {
	e = &emulator{
		sampleBuffer: bytes.NewBuffer(make([]byte, 0)),
		samples:      make([]byte, 4096),
	}
	e.c = core.NewGB(e.sampleBuffer)
}

//export retro_deinit
func retro_deinit() {
	e = nil
}

//export retro_api_version
func retro_api_version() C.uint {
	return 1
}

//export retro_get_system_info
func retro_get_system_info(info *C.struct_retro_system_info) {
	info.library_name = C.CString("DawnGB")
	info.library_version = C.CString("v1")
	info.need_fullpath = C.bool(false)
	info.valid_extensions = C.CString("gb|gbc")
}

//export retro_get_system_av_info
func retro_get_system_av_info(info *C.struct_retro_system_av_info) {
	info.timing.fps = C.double(60.0)
	info.timing.sample_rate = C.double(32768.0)

	info.geometry.base_width = width
	info.geometry.base_height = height
	info.geometry.max_width = width
	info.geometry.max_height = height
	info.geometry.aspect_ratio = C.float(float64(width) / float64(height))
}

//export retro_set_controller_port_device
func retro_set_controller_port_device(port, device C.uint) {
	// nop
}

//export retro_reset
func retro_reset() {
	e.c.Reset(false)
}

//export retro_run
func retro_run() {
	C.call_input_poll_cb()

	if useBitmasks {
		joypadMask := uint(C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.RETRO_DEVICE_ID_JOYPAD_MASK))
		for i := 0; i < len(keymap); i++ {
			pressed := (joypadMask>>keymap[i])&1 == 1
			e.c.SetKeyInput(keymapNames[keymap[i]], pressed)
		}
	} else {
		for i := 0; i < len(keymap); i++ {
			pressed := C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.uint(keymap[i])) != 0
			e.c.SetKeyInput(keymapNames[keymap[i]], pressed)
		}
	}

	e.c.RunFrame()
	audioBatchCallback()
	renderCheckered()
}

func renderCheckered() {
	screen := e.c.Screen()
	buf := make([]uint16, len(screen))
	for i := 0; i < len(screen); i++ {
		r := uint16(screen[i].R >> 3)
		g := uint16(screen[i].G >> 3)
		b := uint16(screen[i].B >> 3)
		buf[i] = (r << 11) | (g << 6) | b
	}
	C.call_video_cb(unsafe.Pointer(&buf[0]), width, height, width*2)
}

//export retro_serialize_size
func retro_serialize_size() C.size_t {
	return 0
}

//export retro_serialize
func retro_serialize(data unsafe.Pointer, size C.size_t) C.bool {
	return false
}

//export retro_unserialize
func retro_unserialize(data unsafe.Pointer, size C.size_t) C.bool {
	return false
}

//export retro_load_game
func retro_load_game(info *C.struct_retro_game_info) C.bool {
	fmt := C.RETRO_PIXEL_FORMAT_RGB565
	C.call_environ_cb(RETRO_ENVIRONMENT_SET_PIXEL_FORMAT, unsafe.Pointer(&fmt))
	useBitmasks = bool(C.call_environ_cb(RETRO_ENVIRONMENT_GET_INPUT_BITMASKS, nil))

	romPath := C.GoString(info.path)
	rom, _ := os.ReadFile(romPath)
	e.c.LoadROM(rom)

	return true
}

func audioBatchCallback() {
	if e.c != nil {
		for i := 0; i < len(e.samples); i++ {
			e.samples[i] = 0
		}
		n, err := e.sampleBuffer.Read(e.samples)
		if err == nil && n >= 4 {
			C.call_audio_batch_cb((*C.int16_t)(unsafe.Pointer(&e.samples[0])), C.ulong(n/4))
		}
	}
}

//export retro_load_game_special
func retro_load_game_special(gameType C.uint, info unsafe.Pointer, numInfo C.size_t) bool {
	return false
}

//export retro_unload_game
func retro_unload_game() {
	// TODO
}

//export retro_get_region
func retro_get_region() C.uint {
	return RETRO_REGION_NTSC
}

//export retro_get_memory_data
func retro_get_memory_data(id C.uint) unsafe.Pointer {
	// TODO
	return nil
}

//export retro_get_memory_size
func retro_get_memory_size(id C.uint) C.uint {
	// TODO
	return 0
}

//export retro_cheat_reset
func retro_cheat_reset() {
	// nop
}

//export retro_cheat_set
func retro_cheat_set(index C.uint, enabled C.bool, code unsafe.Pointer) {
	// nop
}

func main() {}
