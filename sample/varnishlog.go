package main

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
func cbfLine(cbd varnishapi.Callbackdata) {
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
}

func cbfVxid() {
	fmt.Printf("\n%s << %s:%s >> %d\n", strings.Repeat("*", int(headline.Level)), tnames[int(headline.Trx_type)], rnames[int(headline.Reason)], headline.Vxid)
	fmt.Print(buf)
	buf = ""
	headline = nil
}

func cbfGroup() {
	fmt.Println(strings.Repeat("-", 100))
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
