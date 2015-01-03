package keccak

import (
	"bytes"
	"encoding/hex"
	"hash"
	"io"
	"testing"
)

var tests = []struct {
	f      func() hash.Hash
	name   string
	text   string
	vector string
}{
	{
		f:      newKeccak256,
		name:   "Keccak-256",
		vector: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
	},
	{
		f:      newKeccak512,
		name:   "Keccak-512",
		vector: "0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e",
	},
	{
		f:      New256,
		name:   "SHA3-256",
		vector: "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
	},
	{
		f:      New512,
		name:   "SHA3-512",
		vector: "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26",
	},
	{
		f:      New256,
		name:   "SHA3-256",
		text:   "The quick brown fox jumps over the lazy dog",
		vector: "69070dda01975c8c120c3aada1b282394e7f032fa9cf32f4cb2259a0897dfc04",
	},
}

func TestHash(t *testing.T) {
	for _, tt := range tests {
		vector, err := hex.DecodeString(tt.vector)
		if err != nil {
			t.Errorf("%s(%q): %s", tt.name, tt.text, err)
			continue
		}
		h := tt.f()
		io.WriteString(h, tt.text)
		sum := h.Sum(nil)
		if !bytes.Equal(sum, vector) {
			t.Errorf("%s(%q): want %x, got %x", tt.name, tt.text, vector, sum)
		}
	}
}

func benchmark(b *testing.B, f func() hash.Hash, size int64) {
	var tmp [Size]byte
	var msg [8192]byte
	b.SetBytes(size)
	h := f()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(msg[:size])
		h.Sum(tmp[:0])
	}
}

// Benchmark the Keccak-f permutation function
func Benchmark256_8(b *testing.B)  { benchmark(b, New256, 8) }
func Benchmark256_1k(b *testing.B) { benchmark(b, New256, 1024) }
func Benchmark256_8k(b *testing.B) { benchmark(b, New256, 8192) }

func Benchmark512_8(b *testing.B)  { benchmark(b, New512, 8) }
func Benchmark512_1k(b *testing.B) { benchmark(b, New512, 1024) }
func Benchmark512_8k(b *testing.B) { benchmark(b, New512, 8192) }
