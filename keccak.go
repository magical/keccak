package keccak

// roundGeneric implements one round of the keccak-f[1600] permutation.
func roundGeneric(a [5][5]uint64) [5][5]uint64 {
	// Theta
	var c [5]uint64
	for x := range a {
		c[x] = a[x][0] ^ a[x][1] ^ a[x][2] ^ a[x][3] ^ a[x][4]
	}
	for x := range a {
		x0, x1 := (x+4)%5, (x+1)%5
		a[x][0] ^= c[x0] ^ rotl(c[x1], 1)
		a[x][1] ^= c[x0] ^ rotl(c[x1], 1)
		a[x][2] ^= c[x0] ^ rotl(c[x1], 1)
		a[x][3] ^= c[x0] ^ rotl(c[x1], 1)
		a[x][4] ^= c[x0] ^ rotl(c[x1], 1)
	}

	// Rho and pi
	var b [5][5]uint64
	for x := range a {
		for y := range a[0] {
			x0 := y
			y0 := (x*2 + y*3) % 5
			b[x0][y0] = rotl(a[x][y], rotc[x][y])
		}
	}

	// Chi
	for y := range a[0] {
		c := [5]uint64{b[0][y], b[1][y], b[2][y], b[3][y], b[4][y]}
		a[0][y] = b[0][y] ^ ^c[1] & c[2]
		a[1][y] = b[1][y] ^ ^c[2] & c[3]
		a[2][y] = b[2][y] ^ ^c[3] & c[4]
		a[3][y] = b[3][y] ^ ^c[4] & c[0]
		a[4][y] = b[4][y] ^ ^c[0] & c[1]
	}

	return a
}


func rotl(a uint64, r uint) uint64 {
	return a<<r | a>>(64-r)
}

