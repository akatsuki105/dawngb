package main

// #include "./libretro.h"
import "C"
import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/akatsuki105/dawngb/core/gb"
	"github.com/akatsuki105/dawngb/core/gb/cartridge"
)

const (
	retroApiVersion = 1
)

const (
	DMG_BIOS = "dmg_boot.bin"
	CGB_BIOS = "cgb_boot.bin"
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

var console *gb.GB
var screen = make([]uint16, 160*144)
var sampleBuffer = bytes.NewBuffer(make([]uint8, 0))
var samples = [4096]uint8{}
var systemDir = "./"
var saveDir = "./"
var romData = []uint8{}

var bios = struct {
	exists bool
	data   []uint8
	isCGB  bool
}{}

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
			systemDir = C.GoString(cStr)
		}
	}

	// check save directory
	{
		cStr := C.CString("")
		ok := bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY, unsafe.Pointer(&cStr)))
		if ok {
			saveDir = C.GoString(cStr)
		}
	}

	// check BIOS
	bios.exists = false
	if systemDir != "./" {
		if data, err := os.ReadFile(filepath.Join(systemDir, CGB_BIOS)); err == nil {
			bios.exists = true
			bios.data = data
			bios.isCGB = true
		} else if data, err := os.ReadFile(filepath.Join(systemDir, DMG_BIOS)); err == nil {
			bios.exists = true
			bios.data = data
			bios.isCGB = false
		}
	}
}

//export retro_deinit
func retro_deinit() {
	retro_unload_game()
	bios.exists = false
	bios.data = nil
}

//export retro_api_version
func retro_api_version() C.uint {
	return retroApiVersion
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
	if console == nil {
		return
	}
	width, height := console.Resolution()
	info.timing.fps = C.double(59.7275)
	info.timing.sample_rate = C.double(32768.0)

	info.geometry.base_width = C.uint(width)
	info.geometry.base_height = C.uint(height)
	info.geometry.max_width = C.uint(width)
	info.geometry.max_height = C.uint(height)
	info.geometry.aspect_ratio = C.float(float64(width) / float64(height))
}

//export retro_set_controller_port_device
func retro_set_controller_port_device(port, device C.uint) {
	// nop
}

//export retro_reset
func retro_reset() {
	console.Reset(false)
}

//export retro_run
func retro_run() {
	C.call_input_poll_cb()

	if useBitmasks {
		joypadMask := uint(C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.RETRO_DEVICE_ID_JOYPAD_MASK))
		for i := 0; i < len(keymap); i++ {
			pressed := (joypadMask>>keymap[i])&1 == 1
			console.SetKeyInput(keymapNames[keymap[i]], pressed)
		}
	} else {
		for i := 0; i < len(keymap); i++ {
			pressed := C.call_input_state_cb(0, C.RETRO_DEVICE_JOYPAD, 0, C.uint(keymap[i])) != 0
			console.SetKeyInput(keymapNames[keymap[i]], pressed)
		}
	}

	update()
	render()
}

func update() {
	if console != nil {
		console.RunFrame()

		for i := 0; i < len(samples); i++ {
			samples[i] = 0
		}
		if console != nil {
			n, err := sampleBuffer.Read(samples[:])
			if err == nil && n >= 4 {
				C.call_audio_batch_cb((*C.int16_t)(unsafe.Pointer(&samples[0])), C.ulong(n/4))
			}
		}
	}
}

func render() {
	buffer := console.Screen()
	for i := 0; i < len(buffer); i++ {
		screen[i] = newRGB565(buffer[i])
	}

	width, height := console.Resolution()
	C.call_video_cb(unsafe.Pointer(&screen[0]), C.uint(width), C.uint(height), C.ulong(width*2))
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
	C.call_environ_cb(C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT, unsafe.Pointer(&fmt))
	useBitmasks = bool(C.call_environ_cb(C.RETRO_ENVIRONMENT_GET_INPUT_BITMASKS, nil))

	romPath := C.GoString(info.path)
	data, err := os.ReadFile(romPath)
	if err != nil {
		romData = nil
		return false
	}
	romData = data

	intro := false
	if bios.exists {
		if bios.isCGB {
			console = gb.New(gb.MODEL_CGB, sampleBuffer)
			console.Load(gb.LOAD_BIOS, bios.data)
			intro = true
		} else {
			ext := filepath.Ext(romPath)
			if ext == ".gbc" {
				console = gb.New(gb.MODEL_CGB, sampleBuffer) // DMGのBIOSしかない場合は、CGBでダイレクトに起動
			} else {
				console = gb.New(gb.MODEL_DMG, sampleBuffer)
				console.Load(gb.LOAD_BIOS, bios.data)
				intro = true
			}
		}
	} else {
		console = gb.New(gb.MODEL_CGB, sampleBuffer)
	}

	if err := console.Load(gb.LOAD_ROM, romData); err != nil {
		return false
	}
	console.Reset(intro)
	clear(screen)
	loadSaveData(romPath, intro)

	return true
}

func loadSaveData(romPath string, intro bool) {
	if saveDir != "" {
		filename := filepath.Base(romPath)                    // "AA/BB/GAME.gbc" -> "GAME.gbc"
		ext := filepath.Ext(filename)                         // "GAME.gbc" -> ".gbc"
		savename := strings.ReplaceAll(filename, ext, ".srm") // "GAME.gbc" -> "GAME.srm"
		data, err := os.ReadFile(filepath.Join(saveDir, savename))
		if err == nil {
			console.Load(gb.LOAD_SAVE, data)
			console.Reset(intro)
		}
	}
}

//export retro_load_game_special
func retro_load_game_special(gameType C.uint, info unsafe.Pointer, numInfo C.size_t) bool {
	return false
}

//export retro_unload_game
func retro_unload_game() {
	console = nil
	clear(screen)
}

//export retro_get_region
func retro_get_region() C.uint { return C.RETRO_REGION_NTSC }

//export retro_get_memory_data
func retro_get_memory_data(id C.uint) unsafe.Pointer {
	if console != nil {
		switch id {
		case C.RETRO_MEMORY_SAVE_RAM:
			data, err := console.Dump(gb.DUMP_SAVE)
			if err == nil {
				return unsafe.Pointer(C.CBytes(data))
			}
		}
	}
	return nil
}

//export retro_get_memory_size
func retro_get_memory_size(id C.uint) C.uint {
	if console != nil {
		switch id {
		case C.RETRO_MEMORY_SAVE_RAM:
			if len(romData) >= 0x150 {
				ramSize, ok := cartridge.RAM_SIZES[romData[0x149]]
				if ok {
					return C.uint(ramSize)
				}
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
