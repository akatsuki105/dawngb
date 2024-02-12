NAME := dawngb

ifeq ($(OS),Windows_NT)
EXE := .exe
else
EXE :=
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
GODOT?=$(shell which godot)

.PHONY: build wasm goenv profile godot libretro clean

build:
	go build -o build/$(NAME)$(EXE) ./src/ebi

wasm:
	env GOOS=js GOARCH=wasm go build -o docs/core.wasm ./src/ebi

goenv:
	go env

profile:
	go build -o build/profile/profile ./src/profile

godot: goenv
	CGO_ENABLED=1 \
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_CFLAGS='-Og -g3 -g -fPIC' \
	CGO_LDFLAGS='-Og -g3 -g' \
	go build -gcflags=all="-N -l" -tags tools -buildmode=c-shared -x -trimpath -o "build/libgodotgo-gb-macos-$(GOARCH).dylib" ./src/godot/main.go
	cp "build/libgodotgo-gb-macos-$(GOARCH).dylib" "/Users/akatsuki/Dev/Godot/GameBoy/lib/libgodotgo-gb-macos-$(GOARCH).dylib"
	cp "build/libgodotgo-gb-macos-$(GOARCH).h" "/Users/akatsuki/Dev/Godot/GameBoy/lib/libgodotgo-gb-macos-$(GOARCH).h"

libretro: goenv
	GOARCH=$(GOARCH) CGO_ENABLED=1 go build -buildmode=c-shared -o build/libretro/$(NAME)_$(GOARCH)_libretro.dylib ./src/libretro/main.go

clean:
	rm -rf build

