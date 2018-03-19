package varnishapi


var _tags = map[string]string {
        "Debug": "",
        "Error": "",
        "CLI": "",
        "SessOpen": "",
        "SessClose": "",
        "BackendOpen": "",  // Change key count at varnish41(4->6)
        "BackendStart": "", // 4.1.3~
        "BackendReuse": "",
        "BackendClose": "",
        "HttpGarbage": "",
        "Backend": "",
        "Length": "",
        "FetchError": "",
        "BogoHeader": "",
        "LostHeader": "",
        "TTL": "",
        "Fetch_Body": "",
        "VCL_acl": "",
        "VCL_call": "",
        "VCL_trace": "",
        "VCL_return": "",
        "ReqStart": "client.ip",
        "Hit": "",
        "HitPass": "",
        "HitMiss": "",
        "ExpBan": "",
        "ExpKill": "",
        "WorkThread": "",
        "ESI_xmlerror": "",
        "Hash": "",  // Change log data type(str->bin)
        "Backend_health": "",
        "VCL_Log": "",
        "VCL_Error": "",
        "Gzip": "",
        "Link": "",
        "Begin": "",
        "End": "",
        "VSL": "",
        "Storage": "",
        "Timestamp": "",
        "ReqAcct": "",
        "ESI_BodyBytes": "",  // Only Varnish40X
        "PipeAcct": "",
        "BereqAcct": "",
        "ReqMethod": "req.method",
        "ReqURL": "req.url",
        "ReqProtocol": "req.proto",
        "ReqStatus": "",
        "ReqReason": "",
        "ReqHeader": "req.http.",
        "ReqUnset": "unset req.http.",
        "ReqLost": "",
        "RespMethod": "",
        "RespURL": "",
        "RespProtocol": "resp.proto",
        "RespStatus": "resp.status",
        "RespReason": "resp.reason",
        "RespHeader": "resp.http.",
        "RespUnset": "unset resp.http.",
        "RespLost": "",
        "BereqMethod": "bereq.method",
        "BereqURL": "bereq.url",
        "BereqProtocol": "bereq.proto",
        "BereqStatus": "",
        "BereqReason": "",
        "BereqHeader": "bereq.http.",
        "BereqUnset": "unset bereq.http.",
        "BereqLost": "",
        "BerespMethod": "",
        "BerespURL": "",
        "BerespProtocol": "beresp.proto",
        "BerespStatus": "beresp.status",
        "BerespReason": "beresp.reason",
        "BerespHeader":   "beresp.http.",
        "BerespUnset":    "unset beresp.http.",
        "BerespLost":     "",
        "ObjMethod":      "",
        "ObjURL":         "",
        "ObjProtocol":    "obj.proto",
        "ObjStatus": "obj.status",
        "ObjReason": "obj.reason",
        "ObjHeader": "obj.http.",
        "ObjUnset":     "unset obj.http.",
        "ObjLost":      "",
        "Proxy":        "",  // Only Varnish41x
        "ProxyGarbage": "",  // Only Varnish41x
        "VfpAcct":      "",  // Only Varnish41x
        "Witness":      "",  // Only Varnish41x
        "H2RxHdr":   "",  // Only Varnish50x
        "H2RxBody":  "",  // Only Varnish50x
        "H2TxHdr":   "",  // Only Varnish50x
        "H2TxBody":  "",  // Only Varnish50x
}
