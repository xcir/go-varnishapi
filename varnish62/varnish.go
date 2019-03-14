package varnishapi

// Copyright (c) 2018 Shohei Tanaka(@xcir)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
// https://github.com/xcir/go-varnishapi

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
#include <signal.h>
#include "vapi/vsm.h"
#include "vapi/vsl.h"
#include "vapi/vsc.h"
#include "vapi/voptget.h"
#include "vdef.h"
#include "vut.h"
#include "miniobj.h"

int _callback(void *vsl, struct VSL_transaction **trans, void *priv);
void _sighandler(int sig);

int _stat_iter(void *priv, struct VSC_point *pt);

struct gva_VSL_RECORD{
  uint32_t n0;
  uint32_t n1;
};


*/
import "C"

import (
	"fmt"
	"unsafe"
)

//log
var VSL_tags []string
var VSL_tags_rev map[string]int
var VSLQ_grouping []string
var VSL_tagflags []uint

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

type Callback_line_f func(cbd Callbackdata)
type Callback_f func()
type Callback_sig_f func(sig int) int

var gva_cb_line Callback_line_f
var gva_cb_vxid Callback_f
var gva_cb_group Callback_f
var gva_cb_sig Callback_sig_f

var VUT *C.struct_VUT

//stat
type GVA_VSC_level_desc struct {
	Name  string
	Label string
	Sdesc string
	Ldesc string
}
type GVA_VSC_point struct {
	Name      string
	Val       uint64
	Ctype     string
	Semantics int
	Format    int
	Sdesc     string
	Ldesc     string
	Level     GVA_VSC_level_desc
}

var stats map[string]GVA_VSC_point
var vsm *C.struct_vsm
var vsc *C.struct_vsc

//------------------------
//log
//------------------------

//export _callback
func _callback(vsl unsafe.Pointer, trans **C.struct_VSL_transaction, priv unsafe.Pointer) C.int {

	sz := unsafe.Sizeof(trans)
	tx := uintptr(unsafe.Pointer(trans))
	var length uint16
	if tx == 0 {
		return 0
	}
	var cbd Callbackdata
	cbexec_v := false
	cbexec_g := false
	for {
		t := ((**C.struct_VSL_transaction)(unsafe.Pointer(tx)))
		if *t == nil {
			break
		}
		cbd.Level = uint16((*t).level)
		cbd.Vxid = uint32((*t).vxid)
		cbd.Vxid_parent = uint32((*t).vxid_parent)
		cbd.Reason = uint((*t).reason)
		cbd.Trx_type = uint((*t)._type)

		cbexec_v = false
		for {
			i := C.VSL_Next((*t).c)
			if i < 0 {
				return C.int(i)
			}
			if i == 0 {
				break
			}
			if C.VSL_Match((*C.struct_VSL_data)(vsl), (*t).c) == 0 {
				continue
			}
			cbexec_v = true
			cbexec_g = true

			rc := (*C.struct_gva_VSL_RECORD)(unsafe.Pointer((*t).c.rec.ptr))
			length = uint16(rc.n0 & 0xffff)
			cbd.Tag = uint8(rc.n0 >> 24)
			cbd.Isbin = (VSL_tagflags[cbd.Tag] & (C.SLT_F_BINARY | C.SLT_F_UNSAFE)) > 0

			if rc.n1&0x40000000 > 0 {
				cbd.Marker = "c"
			} else if rc.n1&0x80000000 > 0 {
				cbd.Marker = "b"
			} else {
				cbd.Marker = "-"
			}

			if cbd.Isbin {
				cbd.Databin = C.GoBytes(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr))+uintptr(8)), C.int(length))
			} else {
				cbd.Datastr = C.GoStringN(((*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer((*t).c.rec.ptr)) + uintptr(8)))), C.int(length-1))
			}
			if gva_cb_line != nil {
				gva_cb_line(cbd)
			}
		}
		if gva_cb_vxid != nil && cbexec_v {
			gva_cb_vxid()
		}
		tx += sz
	}
	if gva_cb_group != nil && cbexec_g {
		gva_cb_group()
	}

	return 0
}

//export _sighandler
func _sighandler(sig C.int) {
	if gva_cb_sig != nil {
		sig = C.int(gva_cb_sig(int(sig)))
	}
	C.VUT_Signaled(VUT, sig)
}

func setArg(opts []string) {
	for i := len(opts) - 1; i >= 0; i-- {
		if opts[i][0] != '-' {
			if i > 0 && opts[i-1][0] == '-' {
				C.VUT_Arg(VUT, C.int(opts[i-1][1]), C.CString(opts[i]))
			}
			i--
			continue
		} else {
			C.VUT_Arg(VUT, C.int(opts[i][1]), C.CString(""))
		}
	}
}

func getVariables() {
	if len(VSL_tags) > 0 {
		return
	}
	VSL_tags = make([]string, len(&C.VSL_tags))
	VSL_tags_rev = make(map[string]int, len(&C.VSL_tags))
	VSLQ_grouping = make([]string, len(&C.VSLQ_grouping))
	VSL_tagflags = make([]uint, len(&C.VSL_tagflags))

	for i := 0; i < len(VSL_tags); i++ {
		VSL_tags[i] = C.GoString((&C.VSL_tags)[i])
		VSL_tags_rev[C.GoString((&C.VSL_tags)[i])] = i
	}
	for i := 0; i < len(VSLQ_grouping); i++ {
		VSLQ_grouping[i] = C.GoString((&C.VSLQ_grouping)[i])
	}
	for i := 0; i < len(VSL_tagflags); i++ {
		VSL_tagflags[i] = uint((&C.VSL_tagflags)[i])
	}

}

func LogInit(opts []string, cb_line Callback_line_f, cb_vxid Callback_f, cb_group Callback_f, cb_sig Callback_sig_f) error {
	getVariables()
	if VUT != nil {
		LogFini()
	}
	t := &C.struct_vopt_spec{}
	VUT = C.VUT_Init(C.CString("VarnishVUTproc"), 0, (**C.char)(unsafe.Pointer(C.CString(""))), t)

	if VUT == nil {
		return fmt.Errorf("fail VUT_Init")
	}

	VUT.dispatch_f = (*C.VSLQ_dispatch_f)(unsafe.Pointer(C._callback))
	if cb_line != nil {
		gva_cb_line = cb_line
	}
	if cb_vxid != nil {
		gva_cb_vxid = cb_vxid
	}
	if cb_group != nil {
		gva_cb_group = cb_group
	}
	if cb_sig != nil {
		gva_cb_sig = cb_sig
	}

	if opts != nil {
		setArg(opts)
	}
	C.VUT_Setup(VUT)
	C.VUT_Signal((*C.VUT_sighandler_f)(unsafe.Pointer(C._sighandler)))
	return nil
}

func LogStop() {
	if VUT == nil {
		return
	}
}

func LogRun() {
	if VUT == nil {
		return
	}
	C.VUT_Main(VUT)
}

func LogFini() {
	if VUT == nil {
		return
	}
	C.VUT_Fini(&VUT)
}

//------------------------
//stat
//------------------------

//export _stat_iter
func _stat_iter(priv unsafe.Pointer, pt *C.struct_VSC_point) C.int {
	stats[C.GoString(pt.name)] = GVA_VSC_point{
		Name:      C.GoString(pt.name),
		Val:       uint64(*pt.ptr),
		Ctype:     C.GoString(pt.ctype),
		Semantics: int(pt.semantics),
		Format:    int(pt.format),
		Sdesc:     C.GoString(pt.sdesc),
		Ldesc:     C.GoString(pt.ldesc),
		Level: GVA_VSC_level_desc{
			Name:  C.GoString(pt.level.name),
			Label: C.GoString(pt.level.label),
			Sdesc: C.GoString(pt.level.sdesc),
			Ldesc: C.GoString(pt.level.ldesc),
		},
	}
	return 0
}

func StatInit() error {
	if vsm != nil {
		StatFini()
	}
	vsm = C.VSM_New()
	vsc = C.VSC_New()
	if C.VSM_Attach(vsm, 2) > 0 {
		err := C.GoString(C.VSM_Error(vsm))
		StatFini()
		return fmt.Errorf("%s", err)
	}
	return nil
}

func StatGet() map[string]GVA_VSC_point {
	if vsc == nil {
		return nil
	}
	stats = make(map[string]GVA_VSC_point)
	C.VSC_Iter(vsc, vsm, (*C.VSC_iter_f)(unsafe.Pointer(C._stat_iter)), nil)
	return stats
}

func StatFini() {
	if vsc == nil {
		return
	}
	C.VSC_Destroy(&vsc, vsm)
	C.VSM_Destroy(&vsm)
}
