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
		rotc[x][y] = (i * (i + 1) / 2) % 64
		x, y = y, (2*x+3*y)%5
	}

	err := tmpl.Execute(os.Stdout, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func count(n int) []int {
	var out = make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
func mul(a, b int) int { return a * b }
func mod(a int) int    { return a % 5 }

func afunc(x, y int) string {
	return fmt.Sprintf("a[%d]", y%5*5+x%5)
}

func bfunc(x, y int) string {
	return fmt.Sprintf("b%d%d", x%5, y%5)
}

func rotcfunc(x, y int) int {
	return rotc[x%5][y%5]
}

var funcs = template.FuncMap{
	"count": count,
	"add":   add,
	"sub":   sub,
	"mul":   mul,
	"mod":   mod,
	"a":     afunc,
	"b":     bfunc,
	"rotc":  rotcfunc,
}

var tmpl = template.Must(template.New("keccak").Funcs(funcs).Parse(`
// Generated from go run gen.go
// DO NOT EDIT

package keccak

// round implements one round of the keccak-f[1600] permutation.
func roundGo(a *[25]uint64) {
	{{ range $x := count 5 }}
		var b{{$x}}0, b{{$x}}1, b{{$x}}2, b{{$x}}3, b{{$x}}4 uint64
	{{ end }}

	// Theta
	var c0, c1, c2, c3, c4 uint64
	{{ range $x := count 5 }}
		c{{$x}} = {{a $x 0}} ^ {{a $x 1}} ^ {{a $x 2}} ^ {{a $x 3}} ^ {{a $x 4}}
	{{ end }}
	var d0, d1, d2, d3, d4 uint64
	{{ range $x := count 5 }}
		{{ $x4 := add $x 4 | mod }}
		{{ $x1 := add $x 1 | mod }}
		d{{$x}} = c{{$x4}} ^ (c{{$x1}}<<1 | c{{$x1}}>>63)
		{{ range $y := count 5 }}
			{{b $x $y}} = {{a $x $y}} ^ d{{$x}}
		{{ end }}
	{{ end }}

	{{ range $y := count 5 }}
		// Rho / Pi
		{{ range $x := count 5 }}
			{{ $x0 := add $x (mul $y 3) }}
			{{ $y0 := $x }}
			{{ $b := b $x0 $y0 }}
			{{ $r := rotc $x0 $y0 }}
			c{{$x}} = {{$b}}<<{{$r}} | {{$b}}>>{{sub 64 $r}}
		{{ end }}
		// Chi
		{{ range $x := count 5 }}
			{{ $x1 := add $x 1 | mod }}
			{{ $x2 := add $x 2 | mod }}
			{{a $x $y}} = c{{$x}} ^ (c{{$x2}} &^ c{{$x1}})
		{{ end }}
	{{ end }}
}
`))
