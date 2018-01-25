package main

import "C"

import(
    "fmt"
    "./varnishapi"
)

func main(){
  fmt.Println(varnishapi.Stat())
}
