package main

import "C"

import(
    "fmt"
    "strings"
    "./varnishapi"
)




var buf string = ""
var headline varnishapi.Callbackdata
var tnames =map[int]string{
  0:"unknown",
  1:"sess",
  2:"req",
  3:"bereq",
  4:"raw",
}
var rnames =map[int]string{
  0:"unknown",
  1:"HTTP/1",
  2:"rxreq",
  3:"esi",
  4:"restart",
  5:"pass",
  6:"fetch",
  7:"bgfetch",
  8:"pipe",
}
func cbfl(cbd varnishapi.Callbackdata) int{
  t:=varnishapi.Tag2Var(cbd.Tag,cbd.Datastr)
  buf+=fmt.Sprintf("lv:%d vxid:%d vxidp:%d reason:%v trx:%v thd:%v tag:%s data:%s bin:%v isbin:%v key:%s\n",cbd.Level,cbd.Vxid,cbd.Vxid_parent,rnames[int(cbd.Reason)],tnames[int(cbd.Trx_type)],cbd.Marker,varnishapi.VSL_tags[cbd.Tag],cbd.Datastr,cbd.Databin,cbd.Isbin,t.Key)
  if headline.Vxid==0{
    headline = cbd
  }
  return 0
}
func cbfv() int{

  fmt.Printf("\n%s << %s >> %d\n",strings.Repeat("*",int(headline.Level)),tnames[int(headline.Trx_type)], headline.Vxid)
  fmt.Print(buf)
  buf=""
  headline.Vxid=0
  return 0
}
func cbfg() int{
  fmt.Println(strings.Repeat("-",100))
  return 0
}

func cbsig(sig int) int{
  fmt.Println("hello")
  return sig
}

func main(){
    
    opts:=[]string{"-c","-g","session"}

    varnishapi.LogInit(opts,cbfl,cbfv,cbfg,cbsig)
    varnishapi.LogRun()
    varnishapi.LogFini()
    fmt.Println("Finish")
}