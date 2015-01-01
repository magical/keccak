package keccak

import "hash"

const Size = 256 / 8

const BlockSize = 1600/8 - Size*2

var round = roundGo

// digest implements hash.Hash
type digest struct {
	a   [5][5]uint64 // a[y][x][z]
	buf [BlockSize]byte
	len int
}

func New() hash.Hash {
	return &digest{}
}

func (d *digest) Size() int      { return Size }
func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Reset() {
	*d = digest{}
}

func (d *digest) Write(b []byte) (int, error) {
	written := len(b)
	for len(b) > 0 {
		n := copy(d.buf[d.len:], b)
		d.len += n
		b = b[n:]
		if d.len == BlockSize {
			d.flush()
		}
	}
	return written, nil
}

func (d *digest) flush() {
	b := d.buf[:]
loop:
	for y := range d.a {
		for x := range d.a[0] {
			if len(b) == 0 {
				break loop
			}
			d.a[y][x] ^= le64dec(b)
			b = b[8:]
		}
	}
	keccakf(&d.a)
	d.len = 0
}

func keccakf(a *[5][5]uint64) {
	for i := 0; i < 24; i++ {
		round(a)
		a[0][0] ^= RC[i]
	}
}

func (d0 *digest) Sum(b []byte) []byte {
	d := *d0
	d.buf[d.len] = 0x01
	for i := d.len + 1; i < BlockSize; i++ {
		d.buf[i] = 0
	}
	d.buf[BlockSize-1] |= 0x80
	d.flush()

	b = le64enc(b, d.a[0][0])
	b = le64enc(b, d.a[0][1])
	b = le64enc(b, d.a[0][2])
	b = le64enc(b, d.a[0][3])
	return b
}

func le64dec(b []byte) uint64 {
	return uint64(b[0])<<0 | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func le64enc(b []byte, x uint64) []byte {
	return append(b, byte(x), byte(x>>8), byte(x>>16), byte(x>>24), byte(x>>32), byte(x>>40), byte(x>>48), byte(x>>56))
}
