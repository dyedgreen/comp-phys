package main

// #include "float.h"
import "C"

import (
	"fmt"
	"math"
)

func main() {
	// Determine smallest possible epsilon on this machine
	// We use c's long double to test if the CPU supports extended precision floating point numbers
	size := int64(C.float_acc())
	fmt.Println("This computers 'long double' has the following epsilon precision:")
	fmt.Printf("2^(-%v)\n", size)
	fmt.Println("Go's 32 and 64 bit floating points have the following precision:")
	// The Nextafter functions work the same way C.float_acc is implemented.
	// See https://golang.org/src/math/nextafter.go?s=917:957#L25
	fmt.Printf("float32 %e\n", math.Nextafter32(1, 2)-1)
	fmt.Printf("float64 %e\n", math.Nextafter(1, 2)-1)
	fmt.Println("Go has no support for extended floating point formats (but does provide arbitrary precision\nvia `math/big`).")
}
