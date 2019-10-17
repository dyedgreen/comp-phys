package main

// #include "float.h"
import "C"

import (
	"fmt"
	"math"
)

func main() {
	// Determine smallest possible float on this machine
	size := int64(C.float_acc())
	fmt.Println("This computers 'long double' has the following smallest number n > 0:")
	fmt.Printf("2^(-%v)\n", size)
	fmt.Println("The go compiler (and spec) tells us that our 32 and 64 floats have the following precision:")
	fmt.Printf("float32 %1.20e\n", math.SmallestNonzeroFloat32)
	fmt.Printf("float64 %1.20e\n", math.SmallestNonzeroFloat64)
	fmt.Println("Go has no support for extended (non standard) floating point formats.")
}
