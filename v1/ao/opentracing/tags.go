// Copyright (C) 2017 Librato, Inc. All rights reserved.

package opentracing

import "github.com/opentracing/opentracing-go/ext"

// Map selected OpenTracing tag constants to SolarWinds Observability analogs
var otSWOMap = map[string]string{
	string(ext.Component): "OTComponent",

	string(ext.PeerService):  "RemoteController",
	string(ext.PeerAddress):  "RemoteURL",
	string(ext.PeerHostname): "RemoteHost",

	string(ext.HTTPUrl):        "URL",
	string(ext.HTTPMethod):     "Method",
	string(ext.HTTPStatusCode): "Status",

	string(ext.DBInstance):  "Database",
	string(ext.DBStatement): "Query",
	string(ext.DBType):      "Flavor",

	"resource.name": "TransactionName",
}

func translateTagName(key string) string {
	if k := otSWOMap[key]; k != "" {
		return k
	}
	return key
}
