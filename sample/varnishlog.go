package main

# Copyright (c) 2018 Shohei Tanaka(@xcir)
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions
# are met:
# 1. Redistributions of source code must retain the above copyright
#    notice, this list of conditions and the following disclaimer.
# 2. Redistributions in binary form must reproduce the above copyright
#    notice, this list of conditions and the following disclaimer in the
#    documentation and/or other materials provided with the distribution.
#
# THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
# ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
# OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
# HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
# LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
# OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
# SUCH DAMAGE.

# https://github.com/xcir/go-varnishapi

import "C"

import (
	"../head"
	"fmt"
	"strings"
)

var buf string = ""
var headline *varnishapi.Callbackdata

var tnames = map[int]string{
	0: "unknown",
	1: "sess",
	2: "req",
	3: "bereq",
	4: "raw",
}

var rnames = map[int]string{
	0: "unknown",
	1: "HTTP/1",
	2: "rxreq",
	3: "esi",
	4: "restart",
	5: "pass",
	6: "fetch",
	7: "bgfetch",
	8: "pipe",
}

func cbfLine(cbd varnishapi.Callbackdata) int {
	t := varnishapi.Tag2Var(cbd.Tag, cbd.Datastr)
	buf += fmt.Sprintf("%s lv:%d vxid:%d vxid_parent:%d tag:%s var:%s typs:%s isbin:%v data:",
		strings.Repeat("-", int(cbd.Level)), cbd.Level, cbd.Vxid, cbd.Vxid_parent, varnishapi.VSL_tags[cbd.Tag], t.Key, cbd.Marker, cbd.Isbin)
	if cbd.Isbin {
		buf += fmt.Sprintln(cbd.Databin)
	} else {
		buf += fmt.Sprintln(cbd.Datastr)
	}
	if headline == nil {
		headline = &cbd
	}
	return 0
}

func cbfVxid() int {
	fmt.Printf("\n%s << %s:%s >> %d\n", strings.Repeat("*", int(headline.Level)), tnames[int(headline.Trx_type)], rnames[int(headline.Reason)], headline.Vxid)
	fmt.Print(buf)
	buf = ""
	headline = nil
	return 0
}

func cbfGroup() int {
	fmt.Println(strings.Repeat("-", 100))
	return 0
}

func cbSignal(sig int) int {
	return sig
}

func main() {
	opts := []string{"-c", "-g", "session"}
	varnishapi.LogInit(opts, cbfLine, cbfVxid, cbfGroup, cbSignal)
	defer varnishapi.LogFini()
	varnishapi.LogRun()
}
