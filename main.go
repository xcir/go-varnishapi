package main

import "C"

import(
    "fmt"
    "./varnishapi"
    "os"
    "os/signal"
    "syscall"
)

func cbfl(cbd varnishapi.Callbackdata) int{
  //fmt.Printf("lv:%d vxid:%d vxidp:%d reason:%d trx:%d thd:%d tag:%s data:%s bin:%v isbin:%v\n",cbd.level,cbd.vxid,cbd.vxid_parent,cbd.reason,cbd.trx_type,cbd.marker,VSL_tags[cbd.tag],cbd.datastr,cbd.databin,cbd.isbin)
  fmt.Println(cbd)
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

func handleCtrlC(c chan os.Signal) {
    sig := <-c
    fmt.Println("\nsignal: ", sig)
    varnishapi.LogStop()
}

func main(){
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go handleCtrlC(c)
    
    opts:=[]string{"-c","-g","request"}

    varnishapi.LogInit(opts,cbfl,cbfv,cbfg)
    varnishapi.LogRun()
    varnishapi.LogFini()
    fmt.Println("Finish")
}
