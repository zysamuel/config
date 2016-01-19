package main

import (
	"fmt"
)

type BgpPeer struct {
	Id        int
	PeerIp    string `json:"peerip"`
	PeerState string
	LocalAs   int
	RemoteAs  int
}

type BgpPeers []BgpPeer

var peers BgpPeers
var gblId int

func init() {
}
func findPeer(id int) BgpPeer {
	for _, peer := range peers {
		if peer.Id == id {
			return peer
		}
	}
	return BgpPeer{}
}

func createPeer(ipaddr string, localAs int, remoteAs int) BgpPeer {
	gblId += 1

	peer := BgpPeer{
		Id:        gblId,
		PeerIp:    ipaddr,
		PeerState: "Established",
		LocalAs:   localAs,
		RemoteAs:  remoteAs}
	peers = append(peers, peer)
	return peer

}

func deletePeer(id int) error {

	for i, peer := range peers {
		if peer.Id == id {
			peers = append(peers[:i], peers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Peer with id of %d to delete", id)
}
