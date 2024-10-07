package keccak

import "math/bits"

type Shake struct {
	digest
	initialState *[25]uint64 // todo: unique.Handle?

	running uint8
}

func NewShake128(N, S []byte) *Shake { return newShake(N, S, 128/8) }
func NewShake256(N, S []byte) *Shake { return newShake(N, S, 256/8) }

func newShake(N, S []byte, sizeBytes int) *Shake {
	s := new(Shake)
	s.digest.size = sizeBytes
	if len(N) == 0 && len(S) == 0 {
		s.digest.dsbyte = 0x1f // 1111 10...
	} else {
		// cSHAKE
		s.digest.dsbyte = 0x04 // 00 10...
		rate := s.digest.BlockSize()
		s.digest.Write(leftEncode(uint64(rate)))       // rate in bytes
		s.digest.Write(leftEncode(uint64(len(N)) * 8)) // length of N in bits
		s.digest.Write(N)
		s.digest.Write(leftEncode(uint64(len(S)) * 8)) // length of S in bits
		s.digest.Write(S)
		if s.len > 0 || s.ulen > 0 {
			s.pad8()
			s.flush()
		}
	}
	//s.Reset()
	return s
}

func (s *Shake) pad8() {
	if s.ulen > 0 {
		for i := int(s.ulen); i < len(s.buf); i++ {
			s.buf[i] = 0
		}
		s.a[s.len] ^= le64dec(s.buf[:])
		s.len += 1
		s.ulen = 0
	}
}

// Shake is only resettable if Reset is called before the first Write or Read.
// (The first call to Reset makes a copy of the initial state which is restored
// on subsequent calls.)
func (s *Shake) Reset() {
	if s.running == 0 {
		if s.initialState == nil {
			s.initialState = new([25]uint64)
			*s.initialState = s.a
		}
	} else {
		if s.initialState == nil {
			panic("keccak: Reset called after Read or Write")
		}
		s.a = *s.initialState
		s.buf = [8]byte{}
		s.ulen = 0
		s.len = 0
		s.running = 0
	}
}

func (s *Shake) Write(p []byte) (int, error) {
	s.running = 1
	return s.digest.Write(p)
}

func (s *Shake) Read(p []byte) (int, error) {
	if s.running < 2 && len(p) > 0 {
		s.running = 2

		var dsword uint64
		if s.ulen == 0 {
			dsword = uint64(s.dsbyte)
		} else {
			s.buf[s.ulen] = s.dsbyte
			for i := int(s.ulen) + 1; i < len(s.buf); i++ {
				s.buf[i] = 0
			}
			dsword = le64dec(s.buf[:])
		}
		s.a[s.len] ^= dsword

		bs := s.BlockSize() / 8
		s.a[bs-1] ^= 0x80 << 56

		s.len = bs
		s.ulen = 0

	}
	return s.digest.read(p)
}

func (d *digest) read(p []byte) (int, error) {
	bs := d.BlockSize() / 8
	size := len(p)

	if d.ulen > 0 {
		n := copy(p, d.buf[d.ulen:])
		p = p[n:]
		d.ulen += int8(n)
		if int(d.ulen) == len(d.buf) {
			d.ulen = 0
			d.len += 1
		}
	}

	for len(p) >= 8 {
		if d.len == bs {
			d.squeeze()
		}
		le64enc(p[:0], d.a[d.len])
		p = p[8:]
		d.len += 1
	}

	if len(p) > 0 {
		le64enc(d.buf[:0], d.a[d.len])
		d.ulen = int8(copy(p, d.buf[:]))
	}

	return size, nil
}

func (d *digest) squeeze() {
	//fmt.Printf("Squeezing\n", d.len)
	keccakf(&d.a)
	d.len = 0
}

func leftEncode(x uint64) []byte {
	var out [9]byte
	be64enc(out[1:], x)
	i := bits.LeadingZeros64(x|1) / 8 // 0..7
	out[i] = byte(8 - i)
	return out[i:]
}

func be64enc(b []byte, x uint64) {
	_ = b[7]
	b[0] = byte(x >> 56)
	b[1] = byte(x >> 48)
	b[2] = byte(x >> 40)
	b[3] = byte(x >> 32)
	b[4] = byte(x >> 24)
	b[5] = byte(x >> 16)
	b[6] = byte(x >> 8)
	b[7] = byte(x)
}
