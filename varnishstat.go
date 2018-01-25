package main

import "C"

import(
    "fmt"
    "./varnishapi"
)

func main(){
  s:=varnishapi.Stat()
  for name :=range s{
    fmt.Printf("%50s %20d\n",name,s[name].Val)
  }
  
}
