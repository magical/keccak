package keccak

import "hash"

const Size = 256 / 8

const BlockSize = 1600/8 - Size*2

func round(a *[25]uint64) { roundGo(a) }

// digest implements hash.Hash
type digest struct {
	a      [25]uint64 // a[y][x][z]
	buf    [8]byte    // buf[0:ulen] holds a partial uint64
	ulen   int8
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
	d.buf = [8]byte{}
	d.len = 0
}

func (d *digest) Write(b []byte) (int, error) {
	written := len(b)
	bs := d.BlockSize() / 8
	// fill buf first, if non-empty
	if d.ulen > 0 {
		n := copy(d.buf[d.ulen:], b)
		b = b[n:]
		d.ulen += int8(n)
		// flush?
		if int(d.ulen) == len(d.buf) {
			d.a[d.len] ^= le64dec(d.buf[:])
			d.len += 1
			d.ulen = 0
			if d.len == bs {
				d.flush()
			}
		}
	}
	// xor 8-byte chunks into the state
	for len(b) >= 8 {
		d.a[d.len] ^= le64dec(b)
		b = b[8:]
		d.len += 1
		if d.len == bs {
			d.flush()
		}
	} // len(b) < 8
	// store any remaining bytes
	if len(b) > 0 {
		d.ulen = int8(copy(d.buf[:], b))
	}
	return written, nil
}

func (d *digest) flush() {
	//fmt.Printf("Flushing with %d bytes\n", d.len*8 + int(d.ulen))
	keccakf(&d.a)
	d.len = 0
}

func keccakf(a *[25]uint64) {
	for i := 0; i < 24; i++ {
		round(a)
		a[0] ^= roundc[i]
	}
}

func (d *digest) clone() *digest {
	d0 := *d
	return &d0
}

func (d *digest) Sum(b []byte) []byte {
	d = d.clone()
	if d.ulen == 0 {
		d.a[d.len] ^= uint64(d.dsbyte)
	} else {
		d.buf[d.ulen] = d.dsbyte
		for i := int(d.ulen) + 1; i < len(d.buf); i++ {
			d.buf[i] = 0
		}
		d.a[d.len] ^= le64dec(d.buf[:])
	}

	bs := d.BlockSize() / 8
	d.a[bs-1] |= 0x80 << 56
	//d.len = bs
	d.flush()

	for i := 0; i < d.size/8; i++ {
		b = le64enc(b, d.a[i])
	}
	return b
}

func le64dec(b []byte) uint64 {
	_ = b[7]
	return uint64(b[0])<<0 | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func le64enc(b []byte, x uint64) []byte {
	return append(b, byte(x), byte(x>>8), byte(x>>16), byte(x>>24), byte(x>>32), byte(x>>40), byte(x>>48), byte(x>>56))
}
