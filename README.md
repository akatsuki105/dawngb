# DawnGB

DawnGB is "cycle-count" accurate GameBoy Color emulator written in Go.

You can play on [web](https://dawngb.vercel.app/)!

## Screenshots

<img width="352" alt="prism" src="https://gyazo.com/b8fb82f1c38d618a5693d754e8466bf0.png" />&nbsp;<img width="352" alt="megaman" src="https://gyazo.com/47ab176b9bde4efc2c86f2574fbf250b.png" />
<img width="352" alt="shantae" src="https://gyazo.com/71477827ab99b7c42908292291ba414b.png" />&nbsp;<img width="352" alt="pokered" src="https://gyazo.com/22e9e6adf186408cc7a2b3c6af630bd1.png" />

## Features

- GB(DMG) and GBC(CGB) support
- MBC1, MBC3, MBC5, MBC30 support
- Libretro support
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

## Accuracy

Keep the code as simple as possible, so synchronization is done at each instruction, and drawing is done at once on HBlank.

So game like "Prehistorik Man", which modifies the PPU registers during mid-frame, may not draw correctly.
