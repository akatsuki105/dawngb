package sm83

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

func cb00(c *SM83) { c.rlc(&c.R.BC.Hi) }

func cb01(c *SM83) { c.rlc(&c.R.BC.Lo) }

func cb02(c *SM83) { c.rlc(&c.R.DE.Hi) }

func cb03(c *SM83) { c.rlc(&c.R.DE.Lo) }

func cb04(c *SM83) { c.rlc(&c.R.HL.Hi) }

func cb05(c *SM83) { c.rlc(&c.R.HL.Lo) }

// rlc (hl)
func cb06(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	val = (val << 1) | (val >> 7)
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (val == 0), false, false, util.Bit(val, 0)
}

func cb07(c *SM83) { c.rlc(&c.R.A) }

func cb08(c *SM83) { c.rrc(&c.R.BC.Hi) }

func cb09(c *SM83) { c.rrc(&c.R.BC.Lo) }

func cb0A(c *SM83) { c.rrc(&c.R.DE.Hi) }

func cb0B(c *SM83) { c.rrc(&c.R.DE.Lo) }

func cb0C(c *SM83) { c.rrc(&c.R.HL.Hi) }

func cb0D(c *SM83) { c.rrc(&c.R.HL.Lo) }

// rrc (hl)
func cb0E(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	val = (val << 7) | (val >> 1)
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (val == 0), false, false, util.Bit(val, 7)
}

func cb0F(c *SM83) { c.rrc(&c.R.A) }

func cb10(c *SM83) { c.rl(&c.R.BC.Hi) }

func cb11(c *SM83) { c.rl(&c.R.BC.Lo) }

func cb12(c *SM83) { c.rl(&c.R.DE.Hi) }

func cb13(c *SM83) { c.rl(&c.R.DE.Lo) }

func cb14(c *SM83) { c.rl(&c.R.HL.Hi) }

func cb15(c *SM83) { c.rl(&c.R.HL.Lo) }

// rl (hl)
func cb16(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	carry := util.Bit(val, 7)
	val = (val << 1) | btou8(c.R.F.c)
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (val == 0), false, false, carry
}

func cb17(c *SM83) { c.rl(&c.R.A) }

func cb18(c *SM83) { c.rr(&c.R.BC.Hi) }

func cb19(c *SM83) { c.rr(&c.R.BC.Lo) }

func cb1A(c *SM83) { c.rr(&c.R.DE.Hi) }

func cb1B(c *SM83) { c.rr(&c.R.DE.Lo) }

func cb1C(c *SM83) { c.rr(&c.R.HL.Hi) }

func cb1D(c *SM83) { c.rr(&c.R.HL.Lo) }

// rr (hl)
func cb1E(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	carry := btou8(c.R.F.c)
	c.R.F.c = util.Bit(val, 0)
	val = (val >> 1) | (carry << 7)
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), false, false
}

func cb1F(c *SM83) { c.rr(&c.R.A) }

func cb20(c *SM83) { c.sla(&c.R.BC.Hi) }

func cb21(c *SM83) { c.sla(&c.R.BC.Lo) }

func cb22(c *SM83) { c.sla(&c.R.DE.Hi) }

func cb23(c *SM83) { c.sla(&c.R.DE.Lo) }

func cb24(c *SM83) { c.sla(&c.R.HL.Hi) }

func cb25(c *SM83) { c.sla(&c.R.HL.Lo) }

// sla (hl)
func cb26(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	c.R.F.c = util.Bit(val, 7)
	val <<= 1
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), false, false
}

func cb27(c *SM83) { c.sla(&c.R.A) }

func cb28(c *SM83) { c.sra(&c.R.BC.Hi) }

func cb29(c *SM83) { c.sra(&c.R.BC.Lo) }

func cb2A(c *SM83) { c.sra(&c.R.DE.Hi) }

func cb2B(c *SM83) { c.sra(&c.R.DE.Lo) }

func cb2C(c *SM83) { c.sra(&c.R.HL.Hi) }

func cb2D(c *SM83) { c.sra(&c.R.HL.Lo) }

// sra (hl)
func cb2E(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	c.R.F.c = util.Bit(val, 0)
	val = uint8(int8(val) >> 1)
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), false, false
}

func cb2F(c *SM83) { c.sra(&c.R.A) }

func cb30(c *SM83) { c.swap(&c.R.BC.Hi) }

func cb31(c *SM83) { c.swap(&c.R.BC.Lo) }

func cb32(c *SM83) { c.swap(&c.R.DE.Hi) }

func cb33(c *SM83) { c.swap(&c.R.DE.Lo) }

func cb34(c *SM83) { c.swap(&c.R.HL.Hi) }

func cb35(c *SM83) { c.swap(&c.R.HL.Lo) }

// swap (hl)
func cb36(c *SM83) {
	addr := c.R.HL.Pack()
	val := c.bus.Read(addr)
	val = (val << 4) | (val >> 4)
	c.bus.Write(addr, val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (val == 0), false, false, false
}

func cb37(c *SM83) { c.swap(&c.R.A) }

func cb38(c *SM83) { c.srl(&c.R.BC.Hi) }

func cb39(c *SM83) { c.srl(&c.R.BC.Lo) }

func cb3A(c *SM83) { c.srl(&c.R.DE.Hi) }

func cb3B(c *SM83) { c.srl(&c.R.DE.Lo) }

func cb3C(c *SM83) { c.srl(&c.R.HL.Hi) }

func cb3D(c *SM83) { c.srl(&c.R.HL.Lo) }

// srl (hl)
func cb3E(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	c.R.F.c = util.Bit(val, 0)
	val >>= 1
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), false, false
}

func cb3F(c *SM83) { c.srl(&c.R.A) }

func cb40(c *SM83) { c.bit(c.R.BC.Hi, 0) }

func cb41(c *SM83) { c.bit(c.R.BC.Lo, 0) }

func cb42(c *SM83) { c.bit(c.R.DE.Hi, 0) }

func cb43(c *SM83) { c.bit(c.R.DE.Lo, 0) }

func cb44(c *SM83) { c.bit(c.R.HL.Hi, 0) }

func cb45(c *SM83) { c.bit(c.R.HL.Lo, 0) }

func cb46(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 0) }

func cb47(c *SM83) { c.bit(c.R.A, 0) }

func cb48(c *SM83) { c.bit(c.R.BC.Hi, 1) }

func cb49(c *SM83) { c.bit(c.R.BC.Lo, 1) }

func cb4A(c *SM83) { c.bit(c.R.DE.Hi, 1) }

func cb4B(c *SM83) { c.bit(c.R.DE.Lo, 1) }

func cb4C(c *SM83) { c.bit(c.R.HL.Hi, 1) }

func cb4D(c *SM83) { c.bit(c.R.HL.Lo, 1) }

func cb4E(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 1) }

func cb4F(c *SM83) { c.bit(c.R.A, 1) }

func cb50(c *SM83) { c.bit(c.R.BC.Hi, 2) }

func cb51(c *SM83) { c.bit(c.R.BC.Lo, 2) }

func cb52(c *SM83) { c.bit(c.R.DE.Hi, 2) }

func cb53(c *SM83) { c.bit(c.R.DE.Lo, 2) }

func cb54(c *SM83) { c.bit(c.R.HL.Hi, 2) }

func cb55(c *SM83) { c.bit(c.R.HL.Lo, 2) }

func cb56(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 2) }

func cb57(c *SM83) { c.bit(c.R.A, 2) }

func cb58(c *SM83) { c.bit(c.R.BC.Hi, 3) }

func cb59(c *SM83) { c.bit(c.R.BC.Lo, 3) }

func cb5A(c *SM83) { c.bit(c.R.DE.Hi, 3) }

func cb5B(c *SM83) { c.bit(c.R.DE.Lo, 3) }

func cb5C(c *SM83) { c.bit(c.R.HL.Hi, 3) }

func cb5D(c *SM83) { c.bit(c.R.HL.Lo, 3) }

func cb5E(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 3) }

func cb5F(c *SM83) { c.bit(c.R.A, 3) }

func cb60(c *SM83) { c.bit(c.R.BC.Hi, 4) }

func cb61(c *SM83) { c.bit(c.R.BC.Lo, 4) }

func cb62(c *SM83) { c.bit(c.R.DE.Hi, 4) }

func cb63(c *SM83) { c.bit(c.R.DE.Lo, 4) }

func cb64(c *SM83) { c.bit(c.R.HL.Hi, 4) }

func cb65(c *SM83) { c.bit(c.R.HL.Lo, 4) }

func cb66(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 4) }

func cb67(c *SM83) { c.bit(c.R.A, 4) }

func cb68(c *SM83) { c.bit(c.R.BC.Hi, 5) }

func cb69(c *SM83) { c.bit(c.R.BC.Lo, 5) }

func cb6A(c *SM83) { c.bit(c.R.DE.Hi, 5) }

func cb6B(c *SM83) { c.bit(c.R.DE.Lo, 5) }

func cb6C(c *SM83) { c.bit(c.R.HL.Hi, 5) }

func cb6D(c *SM83) { c.bit(c.R.HL.Lo, 5) }

func cb6E(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 5) }

func cb6F(c *SM83) { c.bit(c.R.A, 5) }

func cb70(c *SM83) { c.bit(c.R.BC.Hi, 6) }

func cb71(c *SM83) { c.bit(c.R.BC.Lo, 6) }

func cb72(c *SM83) { c.bit(c.R.DE.Hi, 6) }

func cb73(c *SM83) { c.bit(c.R.DE.Lo, 6) }

func cb74(c *SM83) { c.bit(c.R.HL.Hi, 6) }

func cb75(c *SM83) { c.bit(c.R.HL.Lo, 6) }

func cb76(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 6) }

func cb77(c *SM83) { c.bit(c.R.A, 6) }

func cb78(c *SM83) { c.bit(c.R.BC.Hi, 7) }

func cb79(c *SM83) { c.bit(c.R.BC.Lo, 7) }

func cb7A(c *SM83) { c.bit(c.R.DE.Hi, 7) }

func cb7B(c *SM83) { c.bit(c.R.DE.Lo, 7) }

func cb7C(c *SM83) { c.bit(c.R.HL.Hi, 7) }

func cb7D(c *SM83) { c.bit(c.R.HL.Lo, 7) }

func cb7E(c *SM83) { c.bit(c.bus.Read(c.R.HL.Pack()), 7) }

func cb7F(c *SM83) { c.bit(c.R.A, 7) }

func cb80(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 0, false) }

func cb81(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 0, false) }

func cb82(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 0, false) }

func cb83(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 0, false) }

func cb84(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 0, false) }

func cb85(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 0, false) }

func cb86(c *SM83) { c.set_hl(0, false) }

func cb87(c *SM83) { c.R.A = util.SetBit(c.R.A, 0, false) }

func cb88(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 1, false) }

func cb89(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 1, false) }

func cb8A(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 1, false) }

func cb8B(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 1, false) }

func cb8C(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 1, false) }

func cb8D(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 1, false) }

func cb8E(c *SM83) { c.set_hl(1, false) }

func cb8F(c *SM83) { c.R.A = util.SetBit(c.R.A, 1, false) }

func cb90(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 2, false) }

func cb91(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 2, false) }

func cb92(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 2, false) }

func cb93(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 2, false) }

func cb94(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 2, false) }

func cb95(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 2, false) }

func cb96(c *SM83) { c.set_hl(2, false) }

func cb97(c *SM83) { c.R.A = util.SetBit(c.R.A, 2, false) }

func cb98(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 3, false) }

func cb99(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 3, false) }

func cb9A(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 3, false) }

func cb9B(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 3, false) }

func cb9C(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 3, false) }

func cb9D(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 3, false) }

func cb9E(c *SM83) { c.set_hl(3, false) }

func cb9F(c *SM83) { c.R.A = util.SetBit(c.R.A, 3, false) }

func cbA0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 4, false) }

func cbA1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 4, false) }

func cbA2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 4, false) }

func cbA3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 4, false) }

func cbA4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 4, false) }

func cbA5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 4, false) }

func cbA6(c *SM83) { c.set_hl(4, false) }

func cbA7(c *SM83) { c.R.A = util.SetBit(c.R.A, 4, false) }

func cbA8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 5, false) }

func cbA9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 5, false) }

func cbAA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 5, false) }

func cbAB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 5, false) }

func cbAC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 5, false) }

func cbAD(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 5, false) }

func cbAE(c *SM83) { c.set_hl(5, false) }

func cbAF(c *SM83) { c.R.A = util.SetBit(c.R.A, 5, false) }

func cbB0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 6, false) }

func cbB1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 6, false) }

func cbB2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 6, false) }

func cbB3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 6, false) }

func cbB4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 6, false) }

func cbB5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 6, false) }

func cbB6(c *SM83) { c.set_hl(6, false) }

func cbB7(c *SM83) { c.R.A = util.SetBit(c.R.A, 6, false) }

func cbB8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 7, false) }

func cbB9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 7, false) }

func cbBA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 7, false) }

func cbBB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 7, false) }

func cbBC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 7, false) }

func cbBD(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 7, false) }

func cbBE(c *SM83) { c.set_hl(7, false) }

func cbBF(c *SM83) { c.R.A = util.SetBit(c.R.A, 7, false) }

func cbC0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 0, true) }

func cbC1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 0, true) }

func cbC2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 0, true) }

func cbC3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 0, true) }

func cbC4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 0, true) }

func cbC5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 0, true) }

func cbC6(c *SM83) { c.set_hl(0, true) }

func cbC7(c *SM83) { c.R.A = util.SetBit(c.R.A, 0, true) }

func cbC8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 1, true) }

func cbC9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 1, true) }

func cbCA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 1, true) }

func cbCB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 1, true) }

func cbCC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 1, true) }

func cbCD(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 1, true) }

func cbCE(c *SM83) { c.set_hl(1, true) }

func cbCF(c *SM83) { c.R.A = util.SetBit(c.R.A, 1, true) }

func cbD0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 2, true) }

func cbD1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 2, true) }

func cbD2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 2, true) }

func cbD3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 2, true) }

func cbD4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 2, true) }

func cbD5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 2, true) }

func cbD6(c *SM83) { c.set_hl(2, true) }

func cbD7(c *SM83) { c.R.A = util.SetBit(c.R.A, 2, true) }

func cbD8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 3, true) } // set 3, b

func cbD9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 3, true) }

func cbDA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 3, true) }

func cbDB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 3, true) }

func cbDC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 3, true) }

func cbDD(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 3, true) }

func cbDE(c *SM83) { c.set_hl(3, true) }

func cbDF(c *SM83) { c.R.A = util.SetBit(c.R.A, 3, true) }

func cbE0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 4, true) }

func cbE1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 4, true) }

func cbE2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 4, true) }

func cbE3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 4, true) }

func cbE4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 4, true) }

func cbE5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 4, true) }

func cbE6(c *SM83) { c.set_hl(4, true) }

func cbE7(c *SM83) { c.R.A = util.SetBit(c.R.A, 4, true) }

func cbE8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 5, true) }

func cbE9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 5, true) }

func cbEA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 5, true) }

func cbEB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 5, true) }

func cbEC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 5, true) }

func cbED(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 5, true) }

func cbEE(c *SM83) { c.set_hl(5, true) }

func cbEF(c *SM83) { c.R.A = util.SetBit(c.R.A, 5, true) }

func cbF0(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 6, true) }

func cbF1(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 6, true) }

func cbF2(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 6, true) }

func cbF3(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 6, true) }

func cbF4(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 6, true) }

func cbF5(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 6, true) }

func cbF6(c *SM83) { c.set_hl(6, true) }

func cbF7(c *SM83) { c.R.A = util.SetBit(c.R.A, 6, true) }

func cbF8(c *SM83) { c.R.BC.Hi = util.SetBit(c.R.BC.Hi, 7, true) }

func cbF9(c *SM83) { c.R.BC.Lo = util.SetBit(c.R.BC.Lo, 7, true) }

func cbFA(c *SM83) { c.R.DE.Hi = util.SetBit(c.R.DE.Hi, 7, true) }

func cbFB(c *SM83) { c.R.DE.Lo = util.SetBit(c.R.DE.Lo, 7, true) }

func cbFC(c *SM83) { c.R.HL.Hi = util.SetBit(c.R.HL.Hi, 7, true) }

func cbFD(c *SM83) { c.R.HL.Lo = util.SetBit(c.R.HL.Lo, 7, true) }

func cbFE(c *SM83) { c.set_hl(7, true) }

func cbFF(c *SM83) { c.R.A = util.SetBit(c.R.A, 7, true) }
