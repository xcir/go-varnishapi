package main

import "C"

import(
    "fmt"
    "./varnishapi"
)

func cbfl(cbd varnishapi.Callbackdata) int{
  t:=varnishapi.Tag2Var(cbd.Tag,cbd.Datastr)
  fmt.Printf("lv:%d vxid:%d vxidp:%d reason:%d trx:%d thd:%d tag:%s data:%s bin:%v isbin:%v key:%s\n",cbd.Level,cbd.Vxid,cbd.Vxid_parent,cbd.Reason,cbd.Trx_type,cbd.Marker,varnishapi.VSL_tags[cbd.Tag],cbd.Datastr,cbd.Databin,cbd.Isbin,t.Key)
//  fmt.Println(cbd)
//  fmt.Println(varnishapi.Tag2Var(cbd.Tag,cbd.Datastr))
  return 0
}
func cbfv() int{
  fmt.Println("vxid")
  return 0
}
func cbfg() int{
  fmt.Println("###########################")
  return 0
}

func cbsig(sig int) int{
  fmt.Println("hello")
  return sig
}

func main(){
    
    opts:=[]string{"-c","-g","request"}

    varnishapi.LogInit(opts,cbfl,cbfv,cbfg,cbsig)
    varnishapi.LogRun()
    varnishapi.LogFini()
    fmt.Println("Finish")
}
