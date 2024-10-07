package keccak

import (
	"bytes"
	"encoding/hex"
	"hash"
	"io"
	"testing"
)

var tests = []struct {
	f    func() hash.Hash
	name string
	text string
	hash string
}{
	{
		f:    newKeccak256,
		name: "Keccak-256",
		hash: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
	},
	{
		f:    newKeccak512,
		name: "Keccak-512",
		hash: "0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e",
	},
	{
		f:    New256,
		name: "SHA3-256",
		hash: "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
	},
	{
		f:    New512,
		name: "SHA3-512",
		hash: "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26",
	},
	{
		f:    New256,
		name: "SHA3-256",
		text: "The quick brown fox jumps over the lazy dog",
		hash: "69070dda01975c8c120c3aada1b282394e7f032fa9cf32f4cb2259a0897dfc04",
	},

	{
		f:    New256,
		name: "SHA3-256",
		text: "a",
		hash: "80084bf2fba02475726feb2cab2d8215eab14bc6bdd8bfb2c8151257032ecd8b",
	},
	{
		f:    New256,
		name: "SHA3-256",
		text: "abcdefg",
		hash: "7d55114476dfc6a2fbeaa10e221a8d0f32fc8f2efb69a6e878f4633366917a62",
	},
	{
		f:    New256,
		name: "SHA3-256",
		text: "abcdefgh",
		hash: "3e2020725a38a48eb3bbf75767f03a22c6b3f41f459c831309b06433ec649779",
	},
	{
		f:    New256,
		name: "SHA3-256",
		text: "abcdefghi",
		hash: "f74eb337992307c22bc59eb43e59583a683f3b93077e7f2472508e8c464d2657",
	},
}

func TestHash(t *testing.T) {
	for _, tt := range tests {
		want, err := hex.DecodeString(tt.hash)
		if err != nil {
			t.Errorf("%s(%q): %s", tt.name, tt.text, err)
			continue
		}
		h := tt.f()
		io.WriteString(h, tt.text)
		got := h.Sum(nil)
		if !bytes.Equal(got, want) {
			t.Errorf("%s(%q) = %x, want %x", tt.name, tt.text, got, want)
		}
	}
}

func TestHashSmallWrites(t *testing.T) {
	for _, tt := range tests {
		want, err := hex.DecodeString(tt.hash)
		if err != nil {
			t.Errorf("%s(%q): %s", tt.name, tt.text, err)
			continue
		}
		h := tt.f()
		for i := range []byte(tt.text) {
			io.WriteString(h, tt.text[i:i+1])
		}
		got := h.Sum(nil)
		if !bytes.Equal(got, want) {
			t.Errorf("%s(%q) = %x, want %x", tt.name, tt.text, got, want)
		}
	}
}

func benchmark(b *testing.B, f func() hash.Hash, size int64) {
	var tmp [Size * 2]byte
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
