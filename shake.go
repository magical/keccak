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
		if s.len > 0 {
			s.pad(rate)
			s.flush()
		}
	}
	//s.Reset()
	return s
}

func (s *Shake) pad8() {
	n := -s.len & 7
	for i := 0; i < n; i++ {
		s.buf[s.len+i] = 0
	}
	s.len += n
}

func (s *Shake) pad(rate int) {
	for i := s.len; i < rate && i < len(s.buf); i++ {
		s.buf[i] = 0
	}
	s.len = rate
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
		s.buf = [200]byte{}
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

		s.buf[s.len] = s.dsbyte
		bs := s.BlockSize()
		for i := s.len + 1; i < bs; i++ {
			s.buf[i] = 0
		}
		s.buf[bs-1] |= 0x80

		for i := range s.a {
			if i*8 > bs {
				break
			}
			s.a[i] ^= le64dec(s.buf[i*8:])
		}

		s.len = bs
	}
	return s.digest.read(p)
}

func (d *digest) read(p []byte) (int, error) {
	bs := d.BlockSize()
	size := len(p)
	for len(p) > 0 {
		if d.len == bs {
			d.squeeze(bs)
		}
		n := copy(p, d.buf[:bs])
		d.len += n
		p = p[n:]
	}
	return size, nil
}

func (d *digest) squeeze(bs int) {
	//fmt.Printf("Squeezing\n", d.len)
	keccakf(&d.a)
	b := d.buf[:bs]
	for i := range d.a {
		if len(b) == 0 {
			break
		}
		le64enc(b[:0], d.a[i]) // append
		b = b[8:]
	}
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
