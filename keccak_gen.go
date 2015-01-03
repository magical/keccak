// Generated from go run gen.go
// DO NOT EDIT

package keccak

// round implements one round of the keccak-f[1600] permutation.
func roundGo(a *[25]uint64) {

	var b00, b01, b02, b03, b04 uint64

	var b10, b11, b12, b13, b14 uint64

	var b20, b21, b22, b23, b24 uint64

	var b30, b31, b32, b33, b34 uint64

	var b40, b41, b42, b43, b44 uint64

	// Theta
	var c0, c1, c2, c3, c4 uint64

	c0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]

	c1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]

	c2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]

	c3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]

	c4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]

	var d uint64

	d = c4 ^ (c1<<1 | c1>>63)

	b00 = a[0] ^ d

	b01 = a[5] ^ d

	b02 = a[10] ^ d

	b03 = a[15] ^ d

	b04 = a[20] ^ d

	d = c0 ^ (c2<<1 | c2>>63)

	b10 = a[1] ^ d

	b11 = a[6] ^ d

	b12 = a[11] ^ d

	b13 = a[16] ^ d

	b14 = a[21] ^ d

	d = c1 ^ (c3<<1 | c3>>63)

	b20 = a[2] ^ d

	b21 = a[7] ^ d

	b22 = a[12] ^ d

	b23 = a[17] ^ d

	b24 = a[22] ^ d

	d = c2 ^ (c4<<1 | c4>>63)

	b30 = a[3] ^ d

	b31 = a[8] ^ d

	b32 = a[13] ^ d

	b33 = a[18] ^ d

	b34 = a[23] ^ d

	d = c3 ^ (c0<<1 | c0>>63)

	b40 = a[4] ^ d

	b41 = a[9] ^ d

	b42 = a[14] ^ d

	b43 = a[19] ^ d

	b44 = a[24] ^ d

	// Rho / Pi

	c0 = b00<<0 | b00>>64

	c1 = b11<<44 | b11>>20

	c2 = b22<<43 | b22>>21

	c3 = b33<<21 | b33>>43

	c4 = b44<<14 | b44>>50

	// Chi

	a[0] = c0 ^ (c2 &^ c1)

	a[1] = c1 ^ (c3 &^ c2)

	a[2] = c2 ^ (c4 &^ c3)

	a[3] = c3 ^ (c0 &^ c4)

	a[4] = c4 ^ (c1 &^ c0)

	// Rho / Pi

	c0 = b30<<28 | b30>>36

	c1 = b41<<20 | b41>>44

	c2 = b02<<3 | b02>>61

	c3 = b13<<45 | b13>>19

	c4 = b24<<61 | b24>>3

	// Chi

	a[5] = c0 ^ (c2 &^ c1)

	a[6] = c1 ^ (c3 &^ c2)

	a[7] = c2 ^ (c4 &^ c3)

	a[8] = c3 ^ (c0 &^ c4)

	a[9] = c4 ^ (c1 &^ c0)

	// Rho / Pi

	c0 = b10<<1 | b10>>63

	c1 = b21<<6 | b21>>58

	c2 = b32<<25 | b32>>39

	c3 = b43<<8 | b43>>56

	c4 = b04<<18 | b04>>46

	// Chi

	a[10] = c0 ^ (c2 &^ c1)

	a[11] = c1 ^ (c3 &^ c2)

	a[12] = c2 ^ (c4 &^ c3)

	a[13] = c3 ^ (c0 &^ c4)

	a[14] = c4 ^ (c1 &^ c0)

	// Rho / Pi

	c0 = b40<<27 | b40>>37

	c1 = b01<<36 | b01>>28

	c2 = b12<<10 | b12>>54

	c3 = b23<<15 | b23>>49

	c4 = b34<<56 | b34>>8

	// Chi

	a[15] = c0 ^ (c2 &^ c1)

	a[16] = c1 ^ (c3 &^ c2)

	a[17] = c2 ^ (c4 &^ c3)

	a[18] = c3 ^ (c0 &^ c4)

	a[19] = c4 ^ (c1 &^ c0)

	// Rho / Pi

	c0 = b20<<62 | b20>>2

	c1 = b31<<55 | b31>>9

	c2 = b42<<39 | b42>>25

	c3 = b03<<41 | b03>>23

	c4 = b14<<2 | b14>>62

	// Chi

	a[20] = c0 ^ (c2 &^ c1)

	a[21] = c1 ^ (c3 &^ c2)

	a[22] = c2 ^ (c4 &^ c3)

	a[23] = c3 ^ (c0 &^ c4)

	a[24] = c4 ^ (c1 &^ c0)

}
