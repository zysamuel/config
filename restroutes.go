package main

import (
	"net/http"
)

type ApiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type ApiRoutes []ApiRoute

var routes = ApiRoutes{
	ApiRoute{
		"Index",
		"GET",
		"/",
		Index,
	},
	//ApiRoute{
	//    "BgpPeerCreate",
	//    "GET",
	//    "/bgppeers",
	//    BgpPeerCreate,
	//},
	ApiRoute{
		"BgpPeerShow",
		"GET",
		"/bgppeers/{peerId}",
		ShowBgpPeer,
	},
	ApiRoute{
		"BgpPeerCreate",
		"POST",
		"/bgppeers",
		BgpPeerCreate,
	},
}

