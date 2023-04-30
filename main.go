package main

import (
	"fmt"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	fmt.Printf("hizuru version: %s-%s\n", Version, Revision)
}
