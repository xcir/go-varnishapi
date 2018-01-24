package varnishapi

import(
  "strings"
)

func Tag2Var(tag uint8, data string)GVA_TAGVAR{
  r :=GVA_TAGVAR{}
  var ok bool
  if r.Key, ok= _tags[VSL_tags[tag]]; !ok {
    return r
  }
  t:=strings.SplitN(r.Key," ", 2)
  r.VKey=strings.SplitN(t[len(t)-1],".",2)[0]
  if r.Key==""{
    return r
  }else if(r.Key[len(r.Key)-1:] == "."){
    t = strings.SplitN(data,": ",2)
    r.Key = r.Key + t[0]
    r.Val = ""
    if len(t) > 1{
      r.Val = t[1]
    }
  }else{
    r.Val = data
  }
  return r
}
