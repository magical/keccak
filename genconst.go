// +build ignore

package main

import "fmt"

func main() {
	//Round constants
	var RC [24]uint64
	rc := uint8(1)
	for i := 0; i < 24; i++ {
		for j := 0; j <= 6; j++ {
			RC[i] |= uint64(rc&1) << (1<<uint(j) - 1)
			rc = rc<<1 ^ 0x71&uint8(int8(rc)>>7)
		}
	}
	fmt.Println("package keccak")
	fmt.Printf("var RC = %#016v\n", RC)

	var rot [5][5]uint
	x, y := 1, 0
	for i := 0; i < 24; i++ {
		rot[x][y] = uint((i+1)*(i+2)/2)%64
		x, y = y, (2*x+3*y)%5
	}
	fmt.Printf("var rotc = %#v\n", rot)
}
