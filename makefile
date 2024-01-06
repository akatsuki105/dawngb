NAME := dugb

ifeq ($(OS),Windows_NT)
EXE := .exe
else
EXE :=
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
GODOT?=$(shell which godot)

.PHONY: build wasm goenv godot clean

build:
	go build -o build/$(NAME)$(EXE) ./src/ebi

wasm:
	env GOOS=js GOARCH=wasm go build -o build/$(NAME).wasm ./src/ebi

goenv:
	go env

godot: goenv
	CGO_ENABLED=1 \
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_CFLAGS='-Og -g3 -g -fPIC' \
	CGO_LDFLAGS='-Og -g3 -g' \
	go build -gcflags=all="-N -l" -tags tools -buildmode=c-shared -x -trimpath -o "build/libgodotgo-gb-macos-$(GOARCH).dylib" ./src/godot/main.go

clean:
	rm -rf build

