package keccak

// roundGeneric implements one round of the keccak-f[1600] permutation.
func roundGeneric(a *[5][5]uint64) {
	// Theta
	var c [5]uint64
	for x := range *a {
		c[x] = a[0][x] ^ a[1][x] ^ a[2][x] ^ a[3][x] ^ a[4][x]
	}
	for x := range a[0] {
		x0, x1 := (x+4)%5, (x+1)%5
		a[0][x] ^= c[x0] ^ rotl(c[x1], 1)
		a[1][x] ^= c[x0] ^ rotl(c[x1], 1)
		a[2][x] ^= c[x0] ^ rotl(c[x1], 1)
		a[3][x] ^= c[x0] ^ rotl(c[x1], 1)
		a[4][x] ^= c[x0] ^ rotl(c[x1], 1)
	}

	// Rho and pi
	var b [5][5]uint64
	for y := range *a {
		for x := range a[0] {
			x0 := y
			y0 := (x*2 + y*3) % 5
			b[y0][x0] = rotl(a[y][x], rotc[y][x])
		}
	}

	// Chi
	for y := range a[0] {
		a[y][0] = b[y][0] ^ ^b[y][1]&b[y][2]
		a[y][1] = b[y][1] ^ ^b[y][2]&b[y][3]
		a[y][2] = b[y][2] ^ ^b[y][3]&b[y][4]
		a[y][3] = b[y][3] ^ ^b[y][4]&b[y][0]
		a[y][4] = b[y][4] ^ ^b[y][0]&b[y][1]
	}
}

func rotl(a uint64, r uint) uint64 {
	return a<<(r%64) | a>>(64-(r%64))
}
