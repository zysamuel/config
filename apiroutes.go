package main 

import  ( 
	  "net/http" 
          "github.com/gorilla/mux"
        )
type ApiRoute struct {
    Name         string
    Method   	 string
    Pattern 	 string
    HandlerFunc  http.HandlerFunc 
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

func createNewRestRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
            var handler http.Handler
            handler = Logger(route.HandlerFunc, route.Name)
            router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
    }
    return router
}

