NAME := dawngb

ifeq ($(OS),Windows_NT)
EXE := .exe
else
EXE :=
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
GOOPTIONS ?= -ldflags="-s -w" -trimpath

.PHONY: build wasm profile libretro libretro-deck clean

build:
	go build $(GOOPTIONS) -o build/$(NAME)$(EXE) ./src/ebi

wasm:
	env GOOS=js GOARCH=wasm go build -o docs/core.wasm ./src/ebi

profile:
	go build -o build/profile/profile ./src/profile

libretro:
	GOARCH=$(GOARCH) CGO_ENABLED=1 go build $(GOOPTIONS) -buildmode=c-shared -o ./build/libretro/$(NAME)_libretro.dylib ./src/libretro/main.go

libretro-deck:
	CC=x86_64-unknown-linux-gnu-gcc GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build $(GOOPTIONS) -buildmode=c-shared -o ./build/libretro/$(NAME)_libretro.so ./src/libretro/main.go

clean:
	rm -rf build

