package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/magical/keccak"
)

func main() {
	flag.Parse()
	var sum [keccak.Size]byte
	h := keccak.New256()
	files := flag.Args()
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		h.Reset()
		_, err = io.Copy(h, f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		h.Sum(sum[:0])
		fmt.Printf("%x %s\n", sum, filename)
	}
}
