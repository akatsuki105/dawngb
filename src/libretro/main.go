package main

/*
// libretro.h で RETRO_API がついてる宣言のコメントアウトが必要
// https://github.com/libretro/RetroArch/blob/b443d9974a179ee45c0e5e913b9842c397998193/libretro-common/include/libretro.h
#include "libretro.h"
#include "cfuncs.h"
#include "input.h"
*/
import "C"
import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/akatsuki105/dawngb/core/gb"
)

const AUDIO_BUFFER_SIZE = 4096

const (
	WIDTH  = 160
	HEIGHT = 144
)

const (
	DMG_BIOS = "dmg_boot.bin"
	CGB_BIOS = "cgb_boot.bin"
)

var (
	useBitmasks bool // この機能が有効な場合、入力は(libretro側で規定された)ビットマスクとして一括取得ができる(falseなら、ボタン1つずつ取得する必要がある)
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

type AppState struct {
	GB           *gb.GB
	Screen       []uint16
	ROM          []uint8
	SampleBuffer *bytes.Buffer
	Samples      [AUDIO_BUFFER_SIZE]uint8
	SystemDir    string
	SaveDir      string
	BIOS         struct {
		exists bool
		data   []uint8
		isCGB  bool
	}
	SaveStateBuffer []uint8
	SaveStateSize   int
}

var app AppState = AppState{
	SampleBuffer:    bytes.NewBuffer(make([]uint8, 0, AUDIO_BUFFER_SIZE)),
	SystemDir:       "./",
	SaveDir:         "./",
	SaveStateBuffer: make([]uint8, 0, 32768),
	SaveStateSize:   0,
}

// Environment callback. Gives implementations a way of performing uncommon tasks. Extensible.
//
//export retro_set_environment
func retro_set_environment(cb C.retro_environment_t) { C._retro_set_environment(cb) }

//export retro_set_video_refresh
func retro_set_video_refresh(cb C.retro_video_refresh_t) { C._retro_set_video_refresh(cb) }

//export retro_set_audio_sample
func retro_set_audio_sample(cb C.retro_audio_sample_t) { C._retro_set_audio_sample(cb) }

//export retro_set_audio_sample_batch
func retro_set_audio_sample_batch(cb C.retro_audio_sample_batch_t) {
	C._retro_set_audio_sample_batch(cb)
}

//export retro_set_input_poll
func retro_set_input_poll(cb C.retro_input_poll_t) { C._retro_set_input_poll(cb) }

//export retro_set_input_state
func retro_set_input_state(cb C.retro_input_state_t) { C._retro_set_input_state(cb) }

//export retro_init
func retro_init() {
	// check system directory
	{
		cStr := C.CString("")
		ok := bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY, unsafe.Pointer(&cStr)))
		if ok {
			app.SystemDir = C.GoString(cStr)
		}
	}

	// check save directory
	{
		cStr := C.CString("")
		ok := bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY, unsafe.Pointer(&cStr)))
		if ok {
			app.SaveDir = C.GoString(cStr)
		}
	}

	// Logging
	{
		// ok := bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE, unsafe.Pointer(&C.logging)))
		// if ok {
		// 	C.log_cb = C.logging.log
		// } else {
		// 	// TODO: Fallback to stderr in Go side
		// }
	}

	// Input
	{
		useBitmasks = bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_INPUT_BITMASKS, nil))
		// C.call_environ_cb(C.RETRO_ENVIRONMENT_SET_CONTROLLER_INFO, unsafe.Pointer(&C.ports[0]))
		// C.call_environ_cb(C.RETRO_ENVIRONMENT_SET_INPUT_DESCRIPTORS, unsafe.Pointer(&C.descriptors_1p[0]))
	}

	// check BIOS
	app.BIOS.exists = false
	if app.SystemDir != "./" {
		if data, err := os.ReadFile(filepath.Join(app.SystemDir, CGB_BIOS)); err == nil {
			app.BIOS.exists = true
			app.BIOS.data = data
			app.BIOS.isCGB = true
		} else if data, err := os.ReadFile(filepath.Join(app.SystemDir, DMG_BIOS)); err == nil {
			app.BIOS.exists = true
			app.BIOS.data = data
			app.BIOS.isCGB = false
		}
	}

	app.Screen = make([]uint16, WIDTH*HEIGHT)
}

//export retro_deinit
func retro_deinit() {
	retro_unload_game()
	app.BIOS.exists = false
	app.BIOS.data = nil
	useBitmasks = false
}

//export retro_api_version
func retro_api_version() C.uint { return C.RETRO_API_VERSION }

//export retro_get_system_info
func retro_get_system_info(info *C.struct_retro_system_info) {
	info.library_name = C.CString("DawnGB")
	info.library_version = C.CString("v1")
	info.need_fullpath = C.bool(false)
	info.valid_extensions = C.CString("gb|gbc")
}

//export retro_get_system_av_info
func retro_get_system_av_info(info *C.struct_retro_system_av_info) {
	if app.GB != nil {
		width, height := app.GB.Resolution()
		info.timing.fps = C.double(float64(4*1024*1024) / 70224)
		info.timing.sample_rate = C.double(32768.0)

		info.geometry.base_width, info.geometry.base_height = C.uint(width), C.uint(height)
		info.geometry.max_width, info.geometry.max_height = C.uint(width), C.uint(height)
		info.geometry.aspect_ratio = C.float(float64(width) / float64(height))
	}
}

//export retro_set_controller_port_device
func retro_set_controller_port_device(port, device C.uint) {
	// nop
}

//export retro_reset
func retro_reset() {
	for i := 0; i < len(app.Samples); i++ {
		app.Samples[i] = 0
	}
	app.GB.Reset()
	app.GB.DirectBoot()
}

//export retro_run
func retro_run() {
	if app.GB != nil {
		pollInput()
		update()
		render()
	}
}

func pollInput() {
	C.call_input_poll_cb()

	joypads := uint16(0)
	if useBitmasks {
		joypads = uint16(C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.RETRO_DEVICE_ID_JOYPAD_MASK))
	} else {
		for i := 0; i < (C.RETRO_DEVICE_ID_JOYPAD_R3 + 1); i++ {
			if C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.uint(i)) != 0 {
				joypads |= 1 << i
			}
		}
	}

	for i := 0; i < len(keymap); i++ {
		pressed := (joypads>>keymap[i])&1 == 1
		app.GB.SetKeyInput(keymapNames[keymap[i]], pressed)
	}
}

func update() {
	app.GB.RunFrame()
	n, err := app.SampleBuffer.Read(app.Samples[:])
	if err == nil && n >= 4 {
		C.call_audio_batch_cb((*C.int16_t)(unsafe.Pointer(&app.Samples[0])), C.ulong(n/4))
	}
}

func render() {
	buffer := app.GB.Screen()
	for i := 0; i < len(buffer); i++ {
		app.Screen[i] = newRGB565(buffer[i])
	}

	width, height := app.GB.Resolution()
	C.call_video_cb(unsafe.Pointer(&app.Screen[0]), C.uint(width), C.uint(height), C.ulong(width*2))
}

//export retro_serialize_size
func retro_serialize_size() C.size_t {
	if app.GB != nil {
		return C.size_t(app.SaveStateSize + len(app.Samples))
	}
	return 0
}

//export retro_serialize
func retro_serialize(data unsafe.Pointer, size C.size_t) C.bool {
	if app.GB != nil {
		// save samples
		for i, b := range app.Samples {
			ptr := (*uint8)(unsafe.Add(data, i))
			*ptr = b
		}

		buf := bytes.NewBuffer(app.SaveStateBuffer)
		if ok := app.GB.Serialize(buf); ok {
			for i, b := range buf.Bytes() {
				ptr := (*uint8)(unsafe.Add(data, len(app.Samples)+i))
				*ptr = b
			}
			return true
		}
	}
	return false
}

//export retro_unserialize
func retro_unserialize(data unsafe.Pointer, size C.size_t) C.bool {
	if app.GB != nil {
		// load samples
		for i := 0; i < len(app.Samples); i++ {
			ptr := (*uint8)(unsafe.Add(data, i))
			app.Samples[i] = *ptr
		}

		buf := bytes.NewBuffer(app.SaveStateBuffer)
		for i := 0; i < int(size)-len(app.Samples); i++ {
			ptr := (*uint8)(unsafe.Add(data, len(app.Samples)+i))
			buf.WriteByte(*ptr)
		}
		return C.bool(app.GB.Deserialize(buf))
	}
	return false
}

//export retro_load_game
func retro_load_game(info *C.struct_retro_game_info) C.bool {
	fmt := C.RETRO_PIXEL_FORMAT_RGB565
	C.call_environ_cb(C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT, unsafe.Pointer(&fmt))

	romPath := C.GoString(info.path)
	data, err := os.ReadFile(romPath)
	if err != nil {
		app.ROM = nil
		return false
	}
	app.ROM = data

	intro := false
	if app.BIOS.exists {
		if app.BIOS.isCGB {
			app.GB = gb.New(gb.MODEL_CGB, app.SampleBuffer)
			app.GB.Load(gb.LOAD_BIOS, app.BIOS.data)
			intro = true
		} else {
			ext := filepath.Ext(romPath)
			if ext == ".gbc" {
				app.GB = gb.New(gb.MODEL_CGB, app.SampleBuffer) // DMGのBIOSしかない場合は、CGBでダイレクトに起動
			} else {
				app.GB = gb.New(gb.MODEL_DMG, app.SampleBuffer)
				app.GB.Load(gb.LOAD_BIOS, app.BIOS.data)
				intro = true
			}
		}
	} else {
		app.GB = gb.New(gb.MODEL_CGB, app.SampleBuffer)
	}

	if err := app.GB.Load(gb.LOAD_ROM, app.ROM); err != nil {
		return false
	}
	app.GB.Reset()
	if !intro {
		app.GB.DirectBoot()
	}
	clear(app.Screen[:])
	loadSaveData(romPath, intro)

	var buf bytes.Buffer
	app.GB.Serialize(&buf)
	app.SaveStateSize = buf.Len()

	return true
}

func loadSaveData(romPath string, intro bool) {
	if app.SaveDir != "" {
		filename := filepath.Base(romPath)                    // "AA/BB/GAME.gbc" -> "GAME.gbc"
		ext := filepath.Ext(filename)                         // "GAME.gbc" -> ".gbc"
		savename := strings.ReplaceAll(filename, ext, ".srm") // "GAME.gbc" -> "GAME.srm"
		data, err := os.ReadFile(filepath.Join(app.SaveDir, savename))
		if err == nil {
			app.GB.Load(gb.LOAD_SAVE, data)
			app.GB.Reset()
			if !intro {
				app.GB.DirectBoot()
			}
		}
	}
}

//export retro_load_game_special
func retro_load_game_special(gameType C.uint, info unsafe.Pointer, numInfo C.size_t) bool {
	return false
}

//export retro_unload_game
func retro_unload_game() {
	app.GB = nil
	app.ROM = []uint8{}
	app.SaveStateBuffer = make([]uint8, 0, 32768)
	app.SaveStateSize = 0
	clear(app.Screen[:])
	clear(app.Samples[:])
}

//export retro_get_region
func retro_get_region() C.uint { return C.RETRO_REGION_NTSC }

//export retro_get_memory_data
func retro_get_memory_data(id C.uint) unsafe.Pointer {
	if app.GB != nil {
		switch id {
		case C.RETRO_MEMORY_SAVE_RAM:
			data, err := app.GB.Dump(gb.DUMP_SAVE)
			if err == nil {
				return unsafe.Pointer(C.CBytes(data))
			}
		}
	}
	return nil
}

//export retro_get_memory_size
func retro_get_memory_size(id C.uint) C.uint {
	if app.GB != nil {
		switch id {
		case C.RETRO_MEMORY_SAVE_RAM:
			if len(app.ROM) >= 0x150 {
				return C.uint(app.GB.Cart.RAMSize())
			}
		}
	}

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

// rrrrrggggggbbbbb
func newRGB565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	r5 := uint16((r >> 11) & 0x1F)
	g6 := uint16((g >> 10) & 0x3F)
	b5 := uint16((b >> 11) & 0x1F)
	return ((r5 << 11) | (g6 << 5) | b5)
}
