package main 

import  ("encoding/json"
	"io"
	"io/ioutil"
   "net/http"
   //"net/url"
	)

func Index (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/jsoni;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
             panic(err)
	}
}

func ShowConfigObject (w http.ResponseWriter, r *http.Request) {
	 logger.Println("####  ShowConfigObject called")
}

func ConfigObjectCreate(w http.ResponseWriter, r *http.Request) {
	 logger.Println("#### ConfigObjectCreate called")
	 resource := r.URL
	 logger.Println("### Requested resource is ", resource)
	 /*body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	 if err != nil {
	 	panic(err)
	 }
 	 if err := r.Body.Close(); err != nil {
		panic(err)
	 }*/
	 return 
}
	
func ShowBgpPeer (w http.ResponseWriter, r *http.Request) {
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
	logger.Println("### Object Create is called")
	p  := createPeer(peer.PeerIp, peer.LocalAs, peer.RemoteAs)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
	  panic(err)
	}
}

