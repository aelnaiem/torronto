package torronto

import (
	"fmt"
)

type Peer struct {
	currentState State
	peers        Peers
	host         String
	port         int
}

func (peer Peer) Insert(filename String) int {
	// add the file to the local node

	// divide the file by chunks and push it out
	// to peers
}

func (peer Peer) Query(status Status) int {
	// not sure what this is for...
}

func (peer Peer) Join() int {
	filelist []File  := nil //shouldn't be nil, need to create the list of files here
	//create the message to join 
	joinMessage:= create_message(peer.host, string(peer.port), 1, fileList)
	

	// send out join message to all peers with your file list

	//why do this explicity?
	//if we get a response in time, set status to connected
	//peers.connectPeer(host, port)
	// else if they don't respond in time, set their status to
	//peers.disconnectPeer(host, port)
}

func (peer Peer) Leave() int {
	// push out unique chunks, least replicated first

	// close all sockets
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
