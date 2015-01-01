
package keccak

// roundGeneric implements one round of the keccak-f[1600] permutation.
func roundGeneric(a *[25]uint64) {
	// Theta
	var c [5]uint64
	for x := range c {
		c[x] = a[0+x] ^ a[5+x] ^ a[10+x] ^ a[15+x] ^ a[20+x]
	}
	for x := 0; x < 5; x++ {
		d := c[(x+4)%5] ^ rotl(c[(x+1)%5], 1)
		for y := 0; y < 5; y++ {
			a[y*5+x] ^= d
		}
	}

	// Rho and pi
	var b [5][5]uint64
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			x0 := y
			y0 := (x*2 + y*3) % 5
			b[y0][x0] = rotl(a[y*5+x], rotc[y][x])
		}
	}

	// Chi
	for y := 0; y < 5; y++ {
		a[y*5+0] = b[y][0] ^ ^b[y][1]&b[y][2]
		a[y*5+1] = b[y][1] ^ ^b[y][2]&b[y][3]
		a[y*5+2] = b[y][2] ^ ^b[y][3]&b[y][4]
		a[y*5+3] = b[y][3] ^ ^b[y][4]&b[y][0]
		a[y*5+4] = b[y][4] ^ ^b[y][0]&b[y][1]
	}
}

func rotl(a uint64, r uint) uint64 {
	return a<<(r%64) | a>>(64-(r%64))
}
