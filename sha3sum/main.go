package main

import (
	"io"
	"os"

	"github.com/magical/keccak"
)

const usage = "usage: sha3sum [-h] files..."

var (
	stdout = os.Stdout
	stderr = os.Stderr
)

func main() {
	// Parse flags
	if len(os.Args) <= 1 {
		stderr.WriteString(usage + "\n")
		return
	}
	args := os.Args[1:]
	for _, a := range args {
		switch a {
		case "-h", "--help":
			stdout.WriteString(usage + "\n")
			return
		case "", "-":
		default:
			if a[0] == '-' {
				stderr.WriteString("sha3sum: unrecognized flag " + a + "\n")
				os.Exit(1)
			}
		}
	}

	// compute file hashes
	var sum [keccak.Size]byte
	h := keccak.New256()
	files := args
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			stderr.WriteString("sha3sum: " + err.Error() + "\n")
			continue
		}
		h.Reset()
		_, err = io.Copy(h, f)
		f.Close()
		if err != nil {
			stderr.WriteString("sha3sum: " + err.Error() + "\n")
			continue
		}
		h.Sum(sum[:0])
		stdout.WriteString(toHex(sum[:]) + " " + filename + "\n")
	}
}

func toHex(data []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, len(data)*2)
	for i, v := range data {
		out[i*2] = digits[v>>4]
		out[i*2+1] = digits[v&15]
	}
	return string(out)
}
