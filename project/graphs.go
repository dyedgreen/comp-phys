package main

import (
	"flag"
	"fmt"
)

var fileType *string = flag.String("format", "png", "set type for saving plotted images")

// Generate the nice graphs for the project report, this is slow ...
func genGraphs() {
	// TODO
	fmt.Printf("Generating Graphs (format .%v) ...\n", *fileType)
}
