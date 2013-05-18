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
	// code
}

func (peer Peer) Query(status Status) int {
	// code
}

func (peer Peer) Join() int {
	// join Peers container

	// push local files to other peers

	// pull other files that don't exist locally

}

func (peer Peer) Leave() int {
	// leave Peers container

	// push out unique chunks, least replicated first

	// close all sockets
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
