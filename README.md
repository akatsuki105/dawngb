# DawnGB

DawnGB is GameBoy Color emulator written in Go.

You can play on [web](https://dawngb.vercel.app/)!

## Screenshots

### Desktop

<img width="300" alt="prism" src="https://gyazo.com/82888eedb9501fb6a7c83c2b76f1fe8a.webp" />&nbsp;<img width="300" alt="megaman" src="https://gyazo.com/6a65b22547c7cddeb07a77ad5400afc4.webp" />
<img width="300" alt="shantae" src="https://gyazo.com/d0293d5fc976614a0322f44b3e6c8130.webp" />&nbsp;<img width="300" alt="pokered" src="https://gyazo.com/043aa023624a1da45f6e8487cf33143d.webp" />

### Browser(Web)

<img width="28.5%" height="20%" src="https://gyazo.com/3bf894c527bdd932aab37e0c82f67091.webp" />&nbsp;&nbsp;&nbsp;&nbsp;<img width="62%" src="https://gyazo.com/9e773470f1db70aad0098e6d98187e4f.webp" />

## Features

- GB(DMG) and GBC(CGB) support
- MBC1, MBC3, MBC5, MBC30 support
- Sound(APU) support
- Libretro support(run `make libretro`)
- Multiplatform support
- Work on Browser([here](https://dawngb.vercel.app/))

## Usage

- Desktop: Run `go run ./src/ebi` and drag and drop a ROM file into the window.
- Browser: Visit [here](https://dawngb.vercel.app/).

Key mapping is as follows:

- `A`: X
- `B`: Z
- `Start`: Enter
- `Select`: Backspace
- `↑` `↓` `←` `→`: Arrow keys

## Internal

```sh
.
├── core  # Emulator core
├── src   # Frontend
└── util  # Utility (should be renamed to "internal" in the future)
```

## Accuracy

Keep the code as simple as possible, so synchronization is done at each instruction, and line-rendering is done at once on HBlank.

So game like "Prehistorik Man", which modifies the PPU registers during mid-frame, may not draw correctly.
