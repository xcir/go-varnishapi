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
//    +-(A) length(uint16)
//  [AAAAAAAA][AAAAAAAA][--------][BBBBBBBB]
//                                 +-(B) type(uint8)
//    +-(C) marker
//  [CCDDDDDD][DDDDDDDD][DDDDDDDD][DDDDDDDD]
//      +-(D) identify
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
	uint16_t length;
	uint8_t  _pad;
	uint8_t  tag;
	uint32_t tagflag;
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
/*
type Callbackdata struct {
    level       uint16
    vxid        uint32
    vxid_parent uint32
    reason      uint8
    _type
    transaction_type
    length
    isbin
    datastr
    databin
}

*/
//export _callback
func _callback(vsl unsafe.Pointer, trans **C.struct_VSL_transaction, priv unsafe.Pointer) C.int {

    sz:= unsafe.Sizeof(trans)
    tx:= uintptr(unsafe.Pointer(trans))
    fmt.Println(">>>>>>>>>>>")
    if tx==0 {
        return 0
    }
    
    for {
        t := ((**C.struct_VSL_transaction)(unsafe.Pointer(tx)))
        if *t == nil {
            break
        }
        level   :=(*t).level
        vxid    :=(*t).vxid
        vxidp   :=(*t).vxid_parent
        reason  :=(*t).reason
        trx_type:=(*t)._type
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

            rc    :=(*C.struct_gva_VSL_RECORD)(unsafe.Pointer((*t).c.rec.ptr))
            length:=rc.length
            tag   :=rc.tag
            isbin :=(VSL_tagflags[tag] & C.SLT_F_BINARY) == 1
            thd_type:=0
            
            if       rc.tagflag & 0x40000000 > 0{
                thd_type = 1
            }else if rc.tagflag & 0x80000000 > 0{
                thd_type = 2
            }

            var data string
            if isbin{
                data=C.GoStringN((((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8))))),C.int(length))
            }else{
                data=C.GoStringN((((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8))))),C.int(length -1))
            }
            
            fmt.Printf("lv:%d vxid:%d vxidp:%d reason:%d trx:%d thd:%d tag:%s data:%s isbin:%v\n",level,vxid,vxidp,reason,trx_type,thd_type,VSL_tags[tag],data,isbin)
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
