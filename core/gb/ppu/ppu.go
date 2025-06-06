package ppu

import (
	"image/color"

	"github.com/akatsuki105/dawngb/core/gb/internal"
	"github.com/akatsuki105/dawngb/core/gb/ppu/renderer"
	"github.com/akatsuki105/dawngb/core/gb/ppu/renderer/software"
)

const KB = 1024

const CYCLE = 2

// 便宜的にPPU構造体に入れているが、PPUチップ内にはなくボード上にある
type VRAM struct {
	Data [16 * KB]uint8
	Bank uint8 // 0 or 1; VBK(0xFF4F)
}

type CPU interface {
	Read(addr uint16) uint8
	IRQ(id int)
	HBlank()
	IsCGBMode() bool // CGBモードかどうか
}

// OAM DMA
type DMA struct {
	Active bool
	Src    uint16
	Until  int64
}

// フレームの間に起きたLCDSTAT IRQに関する情報(フレームの始まりにリセットされる)
type LCDStatIRQInfo struct {
	Triggered bool
	Mode      uint8
	Lx, Ly    uint8
}

/*
SoCに組み込まれているため、`/cpu`にある方が正確ではある
また、コードをシンプルにしたいのでスキャンライン単位で描画を行うことにしている(スキャンライン中にSCX,SCYやWX, WYを変更するようなゲームでは正しく描画されない場合がある)
*/
type PPU struct {
	cpu             CPU
	cycles          int64 // 遅れているサイクル数(8.38MHzのマスターサイクル単位)
	screen          [160 * 144]color.NRGBA
	frameCounter    uint64
	lx, ly          int
	r               renderer.Renderer
	RAM             VRAM
	DMA             DMA
	lcdc, stat, lyc uint8
	OAM             [160]uint8
	Palette         [(4 * 8) * 2]uint16 // 4bppの8パレットが BG と OBJ　の1つずつ
	ioreg           [0x30]uint8
	enableLatch     bool // LCDC.7をセットしてPPUを有効にすると、次のフレームから表示が開始される そうじゃないとゴミが表示される
	objCount        uint8
	bgpi, obpi      uint8

	// For debugging
	StatIRQ LCDStatIRQInfo
}

func New(cpu CPU) *PPU {
	p := &PPU{
		cpu: cpu,
	}
	return p
}

func (p *PPU) Reset() {
	p.r = software.New(p.RAM.Data[:], p.Palette[:], p.OAM[:], p.cpu.IsCGBMode)
	p.frameCounter = 0
	p.lx, p.ly = 0, 0
	p.stat = 0x80
	p.RAM.Bank = 0
	p.objCount = 0
	p.DMA.Active, p.DMA.Src, p.DMA.Until = false, 0, 0
	p.bgpi, p.obpi = 0, 0
	clear(p.Palette[:])
}

func (p *PPU) SkipBIOS() {
	p.Write(0xFF40, 0x91) // LCDC
	p.Write(0xFF47, 0xFC) // BGP
	copy(p.Palette[:4], dmgPalette[:])
	copy(p.Palette[32:36], dmgPalette[:])
}

func (p *PPU) Frame() uint64 {
	return p.frameCounter
}

func (p *PPU) Screen() []color.NRGBA {
	return p.screen[:]
}

func (p *PPU) Run(cycles8MHz int64) {
	if p.DMA.Active {
		p.runDMA(cycles8MHz)
	}

	p.cycles += cycles8MHz
	for p.cycles >= 2 { // 1dot = 4MHz
		p.step()
		p.cycles -= 2
	}
}

func (p *PPU) step() {
	if (p.lcdc & (1 << 7)) != 0 {
		if p.ly < 144 {
			switch p.lx {
			case 0:
				if p.ly == 0 {
					p.StatIRQ.Triggered = false
					p.StatIRQ.Mode, p.StatIRQ.Lx, p.StatIRQ.Ly = 0, 0, 0
				}
				p.scanOAM()
			case 80:
				p.drawing()
			case 252 + (int(p.objCount) * 6):
				p.hblank()
			}
		}
		p.lx++
		if p.lx == 456 {
			p.lx = 0
			p.incrementLY()
		}
	}
}

func (p *PPU) incrementLY() {
	p.objCount = 0
	p.ly++
	switch p.ly {
	case 144:
		p.vblank()
	case 154:
		p.ly = 0
		p.enableLatch = false
		p.frameCounter++
	}
	p.compareLYC()
}

func (p *PPU) compareLYC() {
	oldStat := p.stat
	p.stat = internal.SetBit(p.stat, 2, p.ly == int(p.lyc))
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
	}
}

// GBCのBIOSがやる、DMGゲームに対する色付け処理
func (p *PPU) ColorizeDMG() {
	copy(p.Palette[:4], cgbPalette[:])
	copy(p.Palette[32:36], cgbPalette[4:])
	copy(p.Palette[36:40], cgbPalette[4:])
}

func (p *PPU) runDMA(cycles8MHz int64) {
	p.DMA.Until -= cycles8MHz
	if p.DMA.Until <= 0 {
		for i := uint16(0); i < 160; i++ {
			p.Write(0xFE00+i, p.cpu.Read(p.DMA.Src+i))
		}
		p.DMA.Active = false
	}
}

func (p *PPU) TriggerDMA(src uint16, m int64) {
	if !p.DMA.Active {
		p.DMA.Active = true
		p.DMA.Src = src
		p.DMA.Until = 160 * m
	}
}

func statIRQAsserted(stat uint8) bool {
	if ((stat & (1 << 6)) != 0) && ((stat & (1 << 2)) != 0) {
		return true
	}
	switch stat & 0b11 {
	case 0:
		return ((stat & (1 << 3)) != 0)
	case 1:
		return ((stat & (1 << 4)) != 0)
	case 2:
		return ((stat & (1 << 5)) != 0)
	}
	return false
}
