package cpu

import "github.com/akatsuki105/dawngb/util"

var cbTable = [256]opcode{
	/* 0x00 */ cb00, cb01, cb02, cb03, cb04, cb05, cb06, cb07, cb08, cb09, cb0A, cb0B, cb0C, cb0D, cb0E, cb0F,
	/* 0x10 */ cb10, cb11, cb12, cb13, cb14, cb15, cb16, cb17, cb18, cb19, cb1A, cb1B, cb1C, cb1D, cb1E, cb1F,
	/* 0x20 */ cb20, cb21, cb22, cb23, cb24, cb25, cb26, cb27, cb28, cb29, cb2A, cb2B, cb2C, cb2D, cb2E, cb2F,
	/* 0x30 */ cb30, cb31, cb32, cb33, cb34, cb35, cb36, cb37, cb38, cb39, cb3A, cb3B, cb3C, cb3D, cb3E, cb3F,
	/* 0x40 */ cb40, cb41, cb42, cb43, cb44, cb45, cb46, cb47, cb48, cb49, cb4A, cb4B, cb4C, cb4D, cb4E, cb4F,
	/* 0x50 */ cb50, cb51, cb52, cb53, cb54, cb55, cb56, cb57, cb58, cb59, cb5A, cb5B, cb5C, cb5D, cb5E, cb5F,
	/* 0x60 */ cb60, cb61, cb62, cb63, cb64, cb65, cb66, cb67, cb68, cb69, cb6A, cb6B, cb6C, cb6D, cb6E, cb6F,
	/* 0x70 */ cb70, cb71, cb72, cb73, cb74, cb75, cb76, cb77, cb78, cb79, cb7A, cb7B, cb7C, cb7D, cb7E, cb7F,
	/* 0x80 */ cb80, cb81, cb82, cb83, cb84, cb85, cb86, cb87, cb88, cb89, cb8A, cb8B, cb8C, cb8D, cb8E, cb8F,
	/* 0x90 */ cb90, cb91, cb92, cb93, cb94, cb95, cb96, cb97, cb98, cb99, cb9A, cb9B, cb9C, cb9D, cb9E, cb9F,
	/* 0xA0 */ cbA0, cbA1, cbA2, cbA3, cbA4, cbA5, cbA6, cbA7, cbA8, cbA9, cbAA, cbAB, cbAC, cbAD, cbAE, cbAF,
	/* 0xB0 */ cbB0, cbB1, cbB2, cbB3, cbB4, cbB5, cbB6, cbB7, cbB8, cbB9, cbBA, cbBB, cbBC, cbBD, cbBE, cbBF,
	/* 0xC0 */ cbC0, cbC1, cbC2, cbC3, cbC4, cbC5, cbC6, cbC7, cbC8, cbC9, cbCA, cbCB, cbCC, cbCD, cbCE, cbCF,
	/* 0xD0 */ cbD0, cbD1, cbD2, cbD3, cbD4, cbD5, cbD6, cbD7, cbD8, cbD9, cbDA, cbDB, cbDC, cbDD, cbDE, cbDF,
	/* 0xE0 */ cbE0, cbE1, cbE2, cbE3, cbE4, cbE5, cbE6, cbE7, cbE8, cbE9, cbEA, cbEB, cbEC, cbED, cbEE, cbEF,
	/* 0xF0 */ cbF0, cbF1, cbF2, cbF3, cbF4, cbF5, cbF6, cbF7, cbF8, cbF9, cbFA, cbFB, cbFC, cbFD, cbFE, cbFF,
}

var cbCycles = [256]int64{
	/* 0x00 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0x10 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0x20 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0x30 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0x40 */ 2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	/* 0x50 */ 2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	/* 0x60 */ 2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	/* 0x70 */ 2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	/* 0x80 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0x90 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xA0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xB0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xC0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xD0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xE0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	/* 0xF0 */ 2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
}

func cb00(c *Cpu) { c.rlc(&c.r.bc.hi) }

func cb01(c *Cpu) { c.rlc(&c.r.bc.lo) }

func cb02(c *Cpu) { c.rlc(&c.r.de.hi) }

func cb03(c *Cpu) { c.rlc(&c.r.de.lo) }

func cb04(c *Cpu) { c.rlc(&c.r.hl.hi) }

func cb05(c *Cpu) { c.rlc(&c.r.hl.lo) }

// rlc (hl)
func cb06(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	val = (val << 1) | (val >> 7)
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (val == 0), false, false, util.Bit(val, 0)
}

func cb07(c *Cpu) { c.rlc(&c.r.a) }

func cb08(c *Cpu) { c.rrc(&c.r.bc.hi) }

func cb09(c *Cpu) { c.rrc(&c.r.bc.lo) }

func cb0A(c *Cpu) { c.rrc(&c.r.de.hi) }

func cb0B(c *Cpu) { c.rrc(&c.r.de.lo) }

func cb0C(c *Cpu) { c.rrc(&c.r.hl.hi) }

func cb0D(c *Cpu) { c.rrc(&c.r.hl.lo) }

// rrc (hl)
func cb0E(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	val = (val << 7) | (val >> 1)
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (val == 0), false, false, util.Bit(val, 7)
}

func cb0F(c *Cpu) { c.rrc(&c.r.a) }

func cb10(c *Cpu) { c.rl(&c.r.bc.hi) }

func cb11(c *Cpu) { c.rl(&c.r.bc.lo) }

func cb12(c *Cpu) { c.rl(&c.r.de.hi) }

func cb13(c *Cpu) { c.rl(&c.r.de.lo) }

func cb14(c *Cpu) { c.rl(&c.r.hl.hi) }

func cb15(c *Cpu) { c.rl(&c.r.hl.lo) }

// rl (hl)
func cb16(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	carry := util.Bit(val, 7)
	val = (val << 1) | util.Btou8(c.r.f.c)
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (val == 0), false, false, carry
}

func cb17(c *Cpu) { c.rl(&c.r.a) }

func cb18(c *Cpu) { c.rr(&c.r.bc.hi) }

func cb19(c *Cpu) { c.rr(&c.r.bc.lo) }

func cb1A(c *Cpu) { c.rr(&c.r.de.hi) }

func cb1B(c *Cpu) { c.rr(&c.r.de.lo) }

func cb1C(c *Cpu) { c.rr(&c.r.hl.hi) }

func cb1D(c *Cpu) { c.rr(&c.r.hl.lo) }

// rr (hl)
func cb1E(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	carry := util.Btou8(c.r.f.c)
	c.r.f.c = util.Bit(val, 0)
	val = (val >> 1) | (carry << 7)
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), false, false
}

func cb1F(c *Cpu) { c.rr(&c.r.a) }

func cb20(c *Cpu) { c.sla(&c.r.bc.hi) }

func cb21(c *Cpu) { c.sla(&c.r.bc.lo) }

func cb22(c *Cpu) { c.sla(&c.r.de.hi) }

func cb23(c *Cpu) { c.sla(&c.r.de.lo) }

func cb24(c *Cpu) { c.sla(&c.r.hl.hi) }

func cb25(c *Cpu) { c.sla(&c.r.hl.lo) }

// sla (hl)
func cb26(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	c.r.f.c = util.Bit(val, 7)
	val <<= 1
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), false, false
}

func cb27(c *Cpu) { c.sla(&c.r.a) }

func cb28(c *Cpu) { c.sra(&c.r.bc.hi) }

func cb29(c *Cpu) { c.sra(&c.r.bc.lo) }

func cb2A(c *Cpu) { c.sra(&c.r.de.hi) }

func cb2B(c *Cpu) { c.sra(&c.r.de.lo) }

func cb2C(c *Cpu) { c.sra(&c.r.hl.hi) }

func cb2D(c *Cpu) { c.sra(&c.r.hl.lo) }

// sra (hl)
func cb2E(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	c.r.f.c = util.Bit(val, 0)
	val = uint8(int8(val) >> 1)
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), false, false
}

func cb2F(c *Cpu) { c.sra(&c.r.a) }

func cb30(c *Cpu) { c.swap(&c.r.bc.hi) }

func cb31(c *Cpu) { c.swap(&c.r.bc.lo) }

func cb32(c *Cpu) { c.swap(&c.r.de.hi) }

func cb33(c *Cpu) { c.swap(&c.r.de.lo) }

func cb34(c *Cpu) { c.swap(&c.r.hl.hi) }

func cb35(c *Cpu) { c.swap(&c.r.hl.lo) }

// swap (hl)
func cb36(c *Cpu) {
	addr := c.r.hl.pack()
	val := c.m.Read(addr)
	val = (val << 4) | (val >> 4)
	c.m.Write(addr, val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (val == 0), false, false, false
}

func cb37(c *Cpu) { c.swap(&c.r.a) }

func cb38(c *Cpu) { c.srl(&c.r.bc.hi) }

func cb39(c *Cpu) { c.srl(&c.r.bc.lo) }

func cb3A(c *Cpu) { c.srl(&c.r.de.hi) }

func cb3B(c *Cpu) { c.srl(&c.r.de.lo) }

func cb3C(c *Cpu) { c.srl(&c.r.hl.hi) }

func cb3D(c *Cpu) { c.srl(&c.r.hl.lo) }

// srl (hl)
func cb3E(c *Cpu) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	c.r.f.c = util.Bit(val, 0)
	val >>= 1
	c.m.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), false, false
}

func cb3F(c *Cpu) { c.srl(&c.r.a) }

func cb40(c *Cpu) { c.bit(c.r.bc.hi, 0) }

func cb41(c *Cpu) { c.bit(c.r.bc.lo, 0) }

func cb42(c *Cpu) { c.bit(c.r.de.hi, 0) }

func cb43(c *Cpu) { c.bit(c.r.de.lo, 0) }

func cb44(c *Cpu) { c.bit(c.r.hl.hi, 0) }

func cb45(c *Cpu) { c.bit(c.r.hl.lo, 0) }

func cb46(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 0) }

func cb47(c *Cpu) { c.bit(c.r.a, 0) }

func cb48(c *Cpu) { c.bit(c.r.bc.hi, 1) }

func cb49(c *Cpu) { c.bit(c.r.bc.lo, 1) }

func cb4A(c *Cpu) { c.bit(c.r.de.hi, 1) }

func cb4B(c *Cpu) { c.bit(c.r.de.lo, 1) }

func cb4C(c *Cpu) { c.bit(c.r.hl.hi, 1) }

func cb4D(c *Cpu) { c.bit(c.r.hl.lo, 1) }

func cb4E(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 1) }

func cb4F(c *Cpu) { c.bit(c.r.a, 1) }

func cb50(c *Cpu) { c.bit(c.r.bc.hi, 2) }

func cb51(c *Cpu) { c.bit(c.r.bc.lo, 2) }

func cb52(c *Cpu) { c.bit(c.r.de.hi, 2) }

func cb53(c *Cpu) { c.bit(c.r.de.lo, 2) }

func cb54(c *Cpu) { c.bit(c.r.hl.hi, 2) }

func cb55(c *Cpu) { c.bit(c.r.hl.lo, 2) }

func cb56(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 2) }

func cb57(c *Cpu) { c.bit(c.r.a, 2) }

func cb58(c *Cpu) { c.bit(c.r.bc.hi, 3) }

func cb59(c *Cpu) { c.bit(c.r.bc.lo, 3) }

func cb5A(c *Cpu) { c.bit(c.r.de.hi, 3) }

func cb5B(c *Cpu) { c.bit(c.r.de.lo, 3) }

func cb5C(c *Cpu) { c.bit(c.r.hl.hi, 3) }

func cb5D(c *Cpu) { c.bit(c.r.hl.lo, 3) }

func cb5E(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 3) }

func cb5F(c *Cpu) { c.bit(c.r.a, 3) }

func cb60(c *Cpu) { c.bit(c.r.bc.hi, 4) }

func cb61(c *Cpu) { c.bit(c.r.bc.lo, 4) }

func cb62(c *Cpu) { c.bit(c.r.de.hi, 4) }

func cb63(c *Cpu) { c.bit(c.r.de.lo, 4) }

func cb64(c *Cpu) { c.bit(c.r.hl.hi, 4) }

func cb65(c *Cpu) { c.bit(c.r.hl.lo, 4) }

func cb66(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 4) }

func cb67(c *Cpu) { c.bit(c.r.a, 4) }

func cb68(c *Cpu) { c.bit(c.r.bc.hi, 5) }

func cb69(c *Cpu) { c.bit(c.r.bc.lo, 5) }

func cb6A(c *Cpu) { c.bit(c.r.de.hi, 5) }

func cb6B(c *Cpu) { c.bit(c.r.de.lo, 5) }

func cb6C(c *Cpu) { c.bit(c.r.hl.hi, 5) }

func cb6D(c *Cpu) { c.bit(c.r.hl.lo, 5) }

func cb6E(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 5) }

func cb6F(c *Cpu) { c.bit(c.r.a, 5) }

func cb70(c *Cpu) { c.bit(c.r.bc.hi, 6) }

func cb71(c *Cpu) { c.bit(c.r.bc.lo, 6) }

func cb72(c *Cpu) { c.bit(c.r.de.hi, 6) }

func cb73(c *Cpu) { c.bit(c.r.de.lo, 6) }

func cb74(c *Cpu) { c.bit(c.r.hl.hi, 6) }

func cb75(c *Cpu) { c.bit(c.r.hl.lo, 6) }

func cb76(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 6) }

func cb77(c *Cpu) { c.bit(c.r.a, 6) }

func cb78(c *Cpu) { c.bit(c.r.bc.hi, 7) }

func cb79(c *Cpu) { c.bit(c.r.bc.lo, 7) }

func cb7A(c *Cpu) { c.bit(c.r.de.hi, 7) }

func cb7B(c *Cpu) { c.bit(c.r.de.lo, 7) }

func cb7C(c *Cpu) { c.bit(c.r.hl.hi, 7) }

func cb7D(c *Cpu) { c.bit(c.r.hl.lo, 7) }

func cb7E(c *Cpu) { c.bit(c.m.Read(c.r.hl.pack()), 7) }

func cb7F(c *Cpu) { c.bit(c.r.a, 7) }

func cb80(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 0, false) }

func cb81(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 0, false) }

func cb82(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 0, false) }

func cb83(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 0, false) }

func cb84(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 0, false) }

func cb85(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 0, false) }

func cb86(c *Cpu) { c.set_hl(0, false) }

func cb87(c *Cpu) { c.r.a = util.SetBit(c.r.a, 0, false) }

func cb88(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 1, false) }

func cb89(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 1, false) }

func cb8A(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 1, false) }

func cb8B(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 1, false) }

func cb8C(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 1, false) }

func cb8D(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 1, false) }

func cb8E(c *Cpu) { c.set_hl(1, false) }

func cb8F(c *Cpu) { c.r.a = util.SetBit(c.r.a, 1, false) }

func cb90(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 2, false) }

func cb91(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 2, false) }

func cb92(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 2, false) }

func cb93(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 2, false) }

func cb94(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 2, false) }

func cb95(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 2, false) }

func cb96(c *Cpu) { c.set_hl(2, false) }

func cb97(c *Cpu) { c.r.a = util.SetBit(c.r.a, 2, false) }

func cb98(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 3, false) }

func cb99(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 3, false) }

func cb9A(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 3, false) }

func cb9B(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 3, false) }

func cb9C(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 3, false) }

func cb9D(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 3, false) }

func cb9E(c *Cpu) { c.set_hl(3, false) }

func cb9F(c *Cpu) { c.r.a = util.SetBit(c.r.a, 3, false) }

func cbA0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 4, false) }

func cbA1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 4, false) }

func cbA2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 4, false) }

func cbA3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 4, false) }

func cbA4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 4, false) }

func cbA5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 4, false) }

func cbA6(c *Cpu) { c.set_hl(4, false) }

func cbA7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 4, false) }

func cbA8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 5, false) }

func cbA9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 5, false) }

func cbAA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 5, false) }

func cbAB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 5, false) }

func cbAC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 5, false) }

func cbAD(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 5, false) }

func cbAE(c *Cpu) { c.set_hl(5, false) }

func cbAF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 5, false) }

func cbB0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 6, false) }

func cbB1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 6, false) }

func cbB2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 6, false) }

func cbB3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 6, false) }

func cbB4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 6, false) }

func cbB5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 6, false) }

func cbB6(c *Cpu) { c.set_hl(6, false) }

func cbB7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 6, false) }

func cbB8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 7, false) }

func cbB9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 7, false) }

func cbBA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 7, false) }

func cbBB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 7, false) }

func cbBC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 7, false) }

func cbBD(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 7, false) }

func cbBE(c *Cpu) { c.set_hl(7, false) }

func cbBF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 7, false) }

func cbC0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 0, true) }

func cbC1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 0, true) }

func cbC2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 0, true) }

func cbC3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 0, true) }

func cbC4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 0, true) }

func cbC5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 0, true) }

func cbC6(c *Cpu) { c.set_hl(0, true) }

func cbC7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 0, true) }

func cbC8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 1, true) }

func cbC9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 1, true) }

func cbCA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 1, true) }

func cbCB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 1, true) }

func cbCC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 1, true) }

func cbCD(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 1, true) }

func cbCE(c *Cpu) { c.set_hl(1, true) }

func cbCF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 1, true) }

func cbD0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 2, true) }

func cbD1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 2, true) }

func cbD2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 2, true) }

func cbD3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 2, true) }

func cbD4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 2, true) }

func cbD5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 2, true) }

func cbD6(c *Cpu) { c.set_hl(2, true) }

func cbD7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 2, true) }

func cbD8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 3, true) } // set 3, b

func cbD9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 3, true) }

func cbDA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 3, true) }

func cbDB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 3, true) }

func cbDC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 3, true) }

func cbDD(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 3, true) }

func cbDE(c *Cpu) { c.set_hl(3, true) }

func cbDF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 3, true) }

func cbE0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 4, true) }

func cbE1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 4, true) }

func cbE2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 4, true) }

func cbE3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 4, true) }

func cbE4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 4, true) }

func cbE5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 4, true) }

func cbE6(c *Cpu) { c.set_hl(4, true) }

func cbE7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 4, true) }

func cbE8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 5, true) }

func cbE9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 5, true) }

func cbEA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 5, true) }

func cbEB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 5, true) }

func cbEC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 5, true) }

func cbED(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 5, true) }

func cbEE(c *Cpu) { c.set_hl(5, true) }

func cbEF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 5, true) }

func cbF0(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 6, true) }

func cbF1(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 6, true) }

func cbF2(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 6, true) }

func cbF3(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 6, true) }

func cbF4(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 6, true) }

func cbF5(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 6, true) }

func cbF6(c *Cpu) { c.set_hl(6, true) }

func cbF7(c *Cpu) { c.r.a = util.SetBit(c.r.a, 6, true) }

func cbF8(c *Cpu) { c.r.bc.hi = util.SetBit(c.r.bc.hi, 7, true) }

func cbF9(c *Cpu) { c.r.bc.lo = util.SetBit(c.r.bc.lo, 7, true) }

func cbFA(c *Cpu) { c.r.de.hi = util.SetBit(c.r.de.hi, 7, true) }

func cbFB(c *Cpu) { c.r.de.lo = util.SetBit(c.r.de.lo, 7, true) }

func cbFC(c *Cpu) { c.r.hl.hi = util.SetBit(c.r.hl.hi, 7, true) }

func cbFD(c *Cpu) { c.r.hl.lo = util.SetBit(c.r.hl.lo, 7, true) }

func cbFE(c *Cpu) { c.set_hl(7, true) }

func cbFF(c *Cpu) { c.r.a = util.SetBit(c.r.a, 7, true) }
