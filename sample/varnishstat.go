package main

import "C"

import (
	"../head"
	"fmt"
)

func main() {
	varnishapi.StatInit()
	defer varnishapi.StatFini()
	for name, v := range varnishapi.StatGet() {
		fmt.Printf("%50s %20d %s\n", name, v.Val, v.Sdesc)
	}

}
