package varnishapi
/*
#cgo pkg-config: varnishapi
#cgo LDFLAGS: -lvarnishapi -lm

#include <stdarg.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <stdint.h>
#include "vapi/vsm.h"
#include "vapi/vsl.h"
#include "vapi/voptget.h"
#include "vas.h"
#include "vdef.h"
#include "vut.h"
#include "miniobj.h"

int _callback(void *vsl, struct VSL_transaction **trans, void *priv);
void _sighandler(int sig);

//
//
// * (via vsl_int.h)
// *
// * Shared memory log format
// *
// * The log member points to an array of 32bit unsigned integers containing
// * log records.
// *
// * Each logrecord consist of:
// *	[n]		= ((type & 0xff) << 24) | (length & 0xffff) 
// *	[n + 1]		= ((marker & 0x03) << 30) | (identifier & 0x3fffffff)
// *	[n + 2] ... [m]	= content (NUL-terminated)
// *
// * Logrecords are NUL-terminated so that string functions can be run
// * directly on the shmlog data.
// *
// * Notice that the constants in these macros cannot be changed without
// * changing corresponding magic numbers in varnishd/cache/cache_shmlog.c
// 
//
struct gva_VSL_RECORD{
  uint32_t n0;
  uint32_t n1;
};


*/
import "C"

import(
  "unsafe"
  "strings"
)

var VSL_tags      []string
var VSLQ_grouping []string
var VSL_tagflags  []uint

type Callbackdata struct {
  Level       uint16
  Vxid        uint32
  Vxid_parent uint32
  Reason      uint
  Marker      string
  Trx_type    uint
  Tag         uint8
  Isbin       bool
  Datastr     string
  Databin     []byte
}

type GVA_TAGVAR struct{
  Key  string
  Val  string
  VKey string
}

type Callback_line_f func(cbd Callbackdata) int
type Callback_f      func() int
type Callback_sig_f  func(sig int) int

var gva_cb_line  Callback_line_f
var gva_cb_vxid  Callback_f
var gva_cb_group Callback_f
var gva_cb_sig   Callback_sig_f

var VUT *C.struct_VUT

//export _callback
func _callback(vsl unsafe.Pointer, trans **C.struct_VSL_transaction, priv unsafe.Pointer) C.int {

  sz:= unsafe.Sizeof(trans)
  tx:= uintptr(unsafe.Pointer(trans))
  var length uint16
  if tx==0 {
    return 0
  }
  var cbd Callbackdata
  for {
    t := ((**C.struct_VSL_transaction)(unsafe.Pointer(tx)))
    if *t == nil {
      break
    }
    cbd.Level       =uint16((*t).level)
    cbd.Vxid        =uint32((*t).vxid)
    cbd.Vxid_parent =uint32((*t).vxid_parent)
    cbd.Reason      =uint((*t).reason)
    cbd.Trx_type    =uint((*t)._type)
    for {
      i:= C.VSL_Next((*t).c)
      if i < 0{
        return i
      }
      if i == 0{
        break
      }
      if C.VSL_Match((*C.struct_VSL_data)(vsl), (*t).c) == 0 {
        continue
      }

      rc        :=(*C.struct_gva_VSL_RECORD)(unsafe.Pointer((*t).c.rec.ptr))
      length     =uint16(rc.n0 & 0xffff)
      cbd.Tag    =uint8(rc.n0 >> 24)
      cbd.Isbin  =(VSL_tagflags[cbd.Tag] & C.SLT_F_BINARY) == 1
      
      if       rc.n1 & 0x40000000 > 0{
        cbd.Marker = "c"
      }else if rc.n1 & 0x80000000 > 0{
        cbd.Marker = "b"
      }else{
        cbd.Marker = "-"
      }
      
      if cbd.Isbin{
        cbd.Databin=C.GoBytes(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8)), C.int(length))
      }else{
        cbd.Datastr=C.GoStringN((((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8))))), C.int(length -1))
      }
      if gva_cb_line != nil {gva_cb_line(cbd)}
    }
    if gva_cb_vxid != nil {gva_cb_vxid()}
    tx+=sz
  }
  if gva_cb_group != nil {gva_cb_group()}
  

  return 0
}

//export _sighandler
func _sighandler(sig C.int){
  if gva_cb_sig != nil {
    sig = C.int(gva_cb_sig(int(sig)))
  }
  C.VUT_Signaled(VUT, sig)
}

func setArg(opts []string){
  for i:=len(opts) -1; i>= 0; i--{
    if opts[i][0] != '-'{
      if i >0 && opts[i-1][0] == '-'{
        C.VUT_Arg(VUT, C.int(opts[i-1][1]), C.CString(opts[i]))
      }
      i--
      continue
    }else{
      C.VUT_Arg(VUT, C.int(opts[i][1]), C.CString(""))
    }
  }
}

func init(){
  VSL_tags      = make([]string, len(&C.VSL_tags))
  VSLQ_grouping = make([]string, len(&C.VSLQ_grouping))
  VSL_tagflags  = make([]uint,   len(&C.VSL_tagflags))

  for i :=0; i< len(VSL_tags); i++{
    VSL_tags[i] = C.GoString((&C.VSL_tags)[i])
  }
  for i :=0; i< len(VSLQ_grouping); i++{
    VSLQ_grouping[i] = C.GoString((&C.VSLQ_grouping)[i])
  }
  for i :=0; i< len(VSL_tagflags); i++{
    VSL_tagflags[i] = uint((&C.VSL_tagflags)[i])
  }
}

func LogInit(opts []string, cb_line Callback_line_f, cb_vxid Callback_f, cb_group Callback_f, cb_sig Callback_sig_f) int{
  if VUT!=nil {LogFini()}
  t :=&C.struct_vopt_spec{}
  VUT =C.VUT_Init(C.CString("VarnishVUTproc"), 0, (**C.char)(unsafe.Pointer(C.CString(""))), t)

  VUT.dispatch_f = (*C.VSLQ_dispatch_f)(unsafe.Pointer(C._callback))
  if cb_line  != nil {gva_cb_line  = cb_line}
  if cb_vxid  != nil {gva_cb_vxid  = cb_vxid}
  if cb_group != nil {gva_cb_group = cb_group}
  if cb_sig   != nil {gva_cb_sig   = cb_sig}

  if opts != nil {setArg(opts)}
  C.VUT_Setup(VUT)
  C.VUT_Signal((*C.VUT_sighandler_f)(unsafe.Pointer(C._sighandler)));
  return 0
}

func LogStop(){
  VUT.sigint = 1
}

func LogRun(){
  if VUT==nil {return}
  C.VUT_Main(VUT)
}

func LogFini(){
  C.VUT_Fini(&VUT)
  VUT = nil
}

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