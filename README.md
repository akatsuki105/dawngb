# DawnGB

DawnGB is "cycle-count" accurate GameBoy Color emulator written in Go.

You can play on [web](https://dawngb.vercel.app/)!

## Screenshots

<img width="352" alt="prism" src="https://gist.github.com/assets/37920078/c55435ed-833a-47cf-b1d8-321bd3dbce76" />&nbsp;<img width="352" alt="megaman" src="https://gist.github.com/assets/37920078/50395ae9-b12d-4311-b217-1f11b675ea35" />
<img width="352" alt="shantae" src="https://gist.github.com/assets/37920078/c5ff2de5-a714-4b1f-b4f9-044cf27be7fa" />&nbsp;<img width="352" alt="pokered" src="https://gist.github.com/assets/37920078/40754cc3-3da8-4fa1-a823-87d5635fa450" />

## Features

- GB(DMG) and GBC(CGB) support
- MBC1, MBC3, MBC5, MBC30 support
- Libretro support
- Work on Browser([here](https://dawngb.vercel.app/))

## Usage

- Desktop: `go run ./src/ebi` and drag and drop a ROM file into the window.
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
