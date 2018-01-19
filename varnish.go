package main
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
    "fmt"
    "unsafe"
)

var VSL_tags      []string
var VSLQ_grouping []string
var VSL_tagflags  []uint

type Callbackdata struct {
    level            uint16
    vxid             uint32
    vxid_parent      uint32
    reason           uint8
    marker           uint8
    trx_type         uint8
    tag              uint8
    length           uint16
    isbin            bool
    datastr          string
    databin          []byte
}

//export _callback
func _callback(vsl unsafe.Pointer, trans **C.struct_VSL_transaction, priv unsafe.Pointer) C.int {

    sz:= unsafe.Sizeof(trans)
    tx:= uintptr(unsafe.Pointer(trans))
    fmt.Println(">>>>>>>>>>>")
    if tx==0 {
        return 0
    }
    var cbd Callbackdata
    for {
        t := ((**C.struct_VSL_transaction)(unsafe.Pointer(tx)))
        if *t == nil {
            break
        }
        cbd.level         =uint16((*t).level)
        cbd.vxid          =uint32((*t).vxid)
        cbd.vxid_parent   =uint32((*t).vxid_parent)
        cbd.reason        =uint8((*t).reason)
        cbd.trx_type      =uint8((*t)._type)

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

            rc      :=(*C.struct_gva_VSL_RECORD)(unsafe.Pointer((*t).c.rec.ptr))
            cbd.length =uint16(rc.n0 & 0xffff)
            cbd.tag    =uint8((rc.n0 & 0xff) << 24)
            cbd.isbin  =(VSL_tagflags[cbd.tag] & C.SLT_F_BINARY) == 1
            
            if       rc.n1 & 0x40000000 > 0{
                cbd.marker = 1
            }else if rc.n1 & 0x80000000 > 0{
                cbd.marker = 2
            }else{
                cbd.marker = 0
            }

            
            if cbd.isbin{
                cbd.databin=C.GoBytes(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8)), C.int(cbd.length))
            }else{
                cbd.datastr=C.GoStringN((((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8))))), C.int(cbd.length -1))
            }
            fmt.Printf("lv:%d vxid:%d vxidp:%d reason:%d trx:%d thd:%d tag:%s data:%s bin:%v isbin:%v\n",cbd.level,cbd.vxid,cbd.vxid_parent,cbd.reason,cbd.trx_type,cbd.marker,VSL_tags[cbd.tag],cbd.datastr,cbd.databin,cbd.isbin)
            fmt.Println(cbd)
            
        }
        tx+=sz
    }
    

    return 0
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


func main(){

    fmt.Println("Main")
    

    t:=&C.struct_vopt_spec{}
    vut:=C.VUT_Init(C.CString("VarnishVUTproc"), 0, (**C.char)(unsafe.Pointer(C.CString(""))), t)

    vut.dispatch_f = (*C.VSLQ_dispatch_f)(unsafe.Pointer(C._callback))


    C.VUT_Setup(vut)
    
    
    C.VUT_Main(vut)
    fmt.Println("finish")
    C.VUT_Fini(&vut)
    
}
