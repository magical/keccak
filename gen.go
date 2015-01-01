// +build ignore

package main

import (
	"fmt"
	"os"
	"text/template"
)

var roundc [24]uint64
var rotc [5][5]int

func main() {
	// Round constants
	rc := uint8(1)
	for i := 0; i < 24; i++ {
		for j := 0; j <= 6; j++ {
			roundc[i] |= uint64(rc&1) << (1<<uint(j) - 1)
			rc = rc<<1 ^ 0x71&uint8(int8(rc)>>7)
		}
	}

	// Rotation constants
	x, y := 1, 0
	for i := 1; i < 25; i++ {
		rotc[x][y] = (i*(i+1)/2) % 64
		x, y = y, (2*x+3*y)%5
	}

	var ctx = struct {
		Rotc [5][5]int
	} {
		rotc,
	}
	err := tmpl.Execute(os.Stdout, &ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func count(n int) []int {
	var out = make([]int, n)
	for i := 0; i < n; i ++ {
		out[i] = i
	}
	return out
}

func add(a, b, m int) int {
	return (a+b) % m
}

func sub(a, b int) int {
	return a - b
}

func mul(a, b int) int {
	return a*b
}

var funcs = template.FuncMap{
	"count": count,
	"add": add,
	"sub": sub,
	"mul": mul,
}

var tmpl = template.Must(template.New("keccak").Funcs(funcs).Parse(`
// Generated from go run gen.go
// DO NOT EDIT

package keccak

// round implements one round of the keccak-f[1600] permutation.
func roundGo(a *[5][5]uint64) {
	{{ range $x := count 5 }}
		var a{{$x}}0, a{{$x}}1, a{{$x}}2, a{{$x}}3, a{{$x}}4 uint64
	{{ end }}

	// Theta
	var c0, c1, c2, c3, c4 uint64
	{{ range $y := count 5 }}
		{{ range $x := count 5 }}
			{{ if eq $y 0 }}
				c{{$x}} = a[{{$y}}][{{$x}}]
			{{ else }}
				c{{$x}} ^= a[{{$y}}][{{$x}}]
			{{ end }}
		{{ end }}
	{{ end }}
	{{ range $x := count 5 }}
		{{ $x0 := add $x 4 5 }}
		{{ $x1 := add $x 1 5 }}
		{{ range $y := count 5 }}
			a{{$x}}{{$y}} = a[{{$y}}][{{$x}}] ^ c{{$x0}} ^ (c{{$x1}}<<1 | c{{$x1}}>>63)
		{{ end }}
	{{ end }}

	// Rho and pi
	{{ range $y := count 5 }}
		{{ range $x := count 5 }}
			{{ $x0 := $y }}
			{{ $y0 := add (mul $x 2) (mul $y 3) 5 }}
			{{ $r := index $.Rotc $x $y }}
			var b{{$x0}}{{$y0}} = a{{$x}}{{$y}}<<{{$r}} | a{{$x}}{{$y}}>>{{sub 64 $r}}
		{{ end }}
	{{ end }}

	// Chi / output
	{{ range $y := count 5 }}
		{{ range $x := count 5 }}
			a[{{$y}}][{{$x}}] = b{{$x}}{{$y}} ^ (b{{add $x 2 5}}{{$y}} &^ b{{add $x 1 5}}{{$y}})
		{{ end }}
	{{ end }}
}
`))
