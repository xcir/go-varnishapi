※golangの勉強がてら作ってるものなんでまだいろいろアレです


not yet released

==================
go-varnishapi
==================


------------------------------------
Connect to libvarnish api by cgo
------------------------------------

:Author: Shohei Tanaka(@xcir)
:Date: xxxx-xx-xx
:Version: 62.trunk
:Support Varnish Version: 5.2.x 6.0.x 6.1.x 6.2.x
:Manual section: 3

For Python(ctypes)
===================
See this link.
https://github.com/xcir/python-varnishapi/


Installation
============

head
        ::

                go get github.com/xcir/go-varnishapi/head

For Varnish6.0.x
        ::

                go get github.com/xcir/go-varnishapi/varnish60

For Varnish5.2.x
        ::

                go get github.com/xcir/go-varnishapi/varnish52


---

Versioning
============
[varnish-version].[library-version]

60.1 is v1 for Varnish6.0.x

DESCRIPTION
============
Connect to libvarnish api by cgo


Connect to VSL(VarnishLog)
--------------------------------

LogInit
-------------------

Prototype
        ::

                func LogInit(opts []string, cb_line Callback_line_f, cb_vxid Callback_f, cb_group Callback_f, cb_sig Callback_sig_f) error

Parameter
        ::

                
                opts     []string
                cb_line  Callback_line_f callback function per line
                cb_vxid  Callback_f      callback function per vxid(call per line, if group option set to raw)
                cb_group Callback_f      callback function per group(raw, vxid, request, session)
                cb_sig   Callback_sig_f  callback signal handler

===================== ======== ======== =========== ===========
callbacktype \\ group raw      vxid     request     session
===================== ======== ======== =========== ===========
cb_line               per line per line per line    per line
cb_vxi                per line per vxid per vxid    per vxid
cb_grop               per line per vxid per request per session
===================== ======== ======== =========== ===========

Return value
        ::

                int
                

Description
        ::

                Open VSL(using VUT)
Example
        ::

                func cbfLine(cbd varnishapi.Callbackdata) {
                	fmt.Println(cbd)
                }
                func cbfVxid() {
                	fmt.Println(strings.Repeat("-", 20))
                }
                func cbfGroup() {
                	fmt.Println(strings.Repeat("/", 20))
                }
                func cbSignal(sig int) int {
                	//Ignore all signal, if set return 0.
                	return sig
                }
                func main() {
                	opts := []string{"-c", "-g", "session"}
                	varnishapi.LogInit(opts, cbfLine, cbfVxid, cbfGroup, cbSignal)
                	defer varnishapi.LogFini()
                	varnishapi.LogRun()
                }

LogStop
-------------------

Prototype
        ::

                func LogStop()

Parameter
        ::

                
                n/a

Return value
        ::

                n/a
                

Description
        ::

                Stop VUT loop
Example
        ::

                XXXXX

LogRun
-------------------

Prototype
        ::

                func LogRun()

Parameter
        ::

                
                n/a

Return value
        ::

                n/a
                

Description
        ::

                Attach to VSL
Example
        ::

                XXXXX


LogFini
-------------------

Prototype
        ::

                func LogFini()

Parameter
        ::

                
                n/a

Return value
        ::

                n/a
                

Description
        ::

                Finish VUT
Example
        ::

                XXXXX



Connect to VSC(VarnishStat)
--------------------------------

StatInit
-------------------

Prototype
        ::

                func StatInit()error

Parameter
        ::

                
                n/a

Return value
        ::

                error
                

Description
        ::

                VSC initialize
Example
        ::

                XXXXX

StatGet
-------------------

Prototype
        ::

                func StatGet()map[string]GVA_VSC_point

Parameter
        ::

                
                n/a

Return value
        ::

                map[string]GVA_VSC_point
                

Description
        ::

                Get VSC values.
Example
        ::

                XXXXX

StatFini
-------------------

Prototype
        ::

                func StatFini()

Parameter
        ::

                
                n/a

Return value
        ::

                n/a
                

Description
        ::

                Finish VSC
Example
        ::

                XXXXX


COPYRIGHT
===========

go-varnishapi

* Copyright (c) 2018 Shohei Tanaka(@xcir)




