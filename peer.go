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
	// send out join message to all peers with your file list
	// if they don't respond in time, set their status to
	// disconnected

	peers.connectPeer(host, port)
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
