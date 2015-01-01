package keccak

import "hash"

const Size = 256 / 8

const BlockSize = 1600/8 - Size*2

var round = roundGo

// digest implements hash.Hash
type digest struct {
	a      [25]uint64 // a[y][x][z]
	buf    [200]byte
	dsbyte byte
	len    int
	size   int
}

func New256() hash.Hash { return &digest{size: 256 / 8, dsbyte: 0x06} }
func New512() hash.Hash { return &digest{size: 512 / 8, dsbyte: 0x06} }

func newKeccak256() hash.Hash { return &digest{size: 256 / 8, dsbyte: 0x01} }
func newKeccak512() hash.Hash { return &digest{size: 512 / 8, dsbyte: 0x01} }

func (d *digest) Size() int      { return d.size }
func (d *digest) BlockSize() int { return 200 - d.size*2 }

func (d *digest) Reset() {
	//fmt.Println("resetting")
	d.a = [25]uint64{}
	d.buf = [200]byte{}
	d.len = 0
}

func (d *digest) Write(b []byte) (int, error) {
	written := len(b)
	bs := d.BlockSize()
	for len(b) > 0 {
		n := copy(d.buf[d.len:bs], b)
		d.len += n
		b = b[n:]
		if d.len == bs {
			d.flush()
		}
	}
	return written, nil
}

func (d *digest) flush() {
	//fmt.Printf("Flushing with %d bytes\n", d.len)
	b := d.buf[:d.len]
loop:
	for i := range d.a {
			if len(b) == 0 {
				break loop
			}
			d.a[i] ^= le64dec(b)
			b = b[8:]
	}
	keccakf(&d.a)
	d.len = 0
}

func keccakf(a *[25]uint64) {
	for i := 0; i < 24; i++ {
		round(a)
		a[0] ^= roundc[i]
	}
}

func (d0 *digest) Sum(b []byte) []byte {
	d := *d0
	d.buf[d.len] = d.dsbyte
	bs := d.BlockSize()
	for i := d.len + 1; i < bs; i++ {
		d.buf[i] = 0
	}
	d.buf[bs-1] |= 0x80
	d.len = bs
	d.flush()

	for i := 0; i < d.size/8; i++ {
		b = le64enc(b, d.a[i])
	}
	return b
}

func le64dec(b []byte) uint64 {
	return uint64(b[0])<<0 | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func le64enc(b []byte, x uint64) []byte {
	return append(b, byte(x), byte(x>>8), byte(x>>16), byte(x>>24), byte(x>>32), byte(x>>40), byte(x>>48), byte(x>>56))
}
