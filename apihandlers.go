package main 

import  ("encoding/json"
	"io"
	"io/ioutil"
	// "fmt"
	// "html"
         "net/http"
         //"github.com/gorilla/mux"
	)

func Index (w http.ResponseWriter, r *http.Request) {
	//peers := BgpPeers {
        //		BgpPeer{ PeerIp    : "10.0.1.1", 
	//		PeerState : "Established",
	//		LocalAs   : 500,
	//		RemoteAs  : 500},
	//	BgpPeer{ PeerIp    : "20.20.2.1",
	//		PeerState : "Established",
	//		LocalAs   : 200,
	//		RemoteAs  : 300},
	//		}
	w.Header().Set("Content-type", "application/jsoni;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
             panic(err)
	}
}

func ShowBgpPeer (w http.ResponseWriter, r *http.Request) {
     //peers := BgpPeers {
     //		BgpPeer{ PeerIp    : "10.0.1.1", 
     //                    PeerState : "Established",
     //			 LocalAs   : 500,	
     //			 RemoteAs  : 500},
     //         BgpPeer{ PeerIp    : "20.20.2.1",
     //                    PeerState : "Established",
     //			 LocalAs   : 200,
     //			 RemoteAs  : 300},
     //         }
     json.NewEncoder(w).Encode(peers)
}

func BgpPeerCreate (w http.ResponseWriter, r *http.Request) {
	var peer BgpPeer
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err:= json.Unmarshal(body, &peer); err != nil  {
		w.Header().Set("Content-Type", "application/json;char-set=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
		    panic(err)
		}
	}
	p  := createPeer(peer.PeerIp, peer.LocalAs, peer.RemoteAs)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
	  panic(err)
	}
}

