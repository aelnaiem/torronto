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
	// join Peers container, by choosing n peers to send messages to

	// push local files to other peers

	// pull other files that don't exist locally

}

func (peer Peer) Leave() int {
	// leave Peers container

	// push out unique chunks, least replicated first

	// close all sockets
}

// func newConnection(newPeer Peer) {
// code to replace newest connection with the newPeer
// Peers.addPeer(peer)
// }

// func updateConnection(hostName string, portNumber int) bool {
//	 check if hostName and portNumber in list of receivers
//	if they are, remove them and add a new receiver, else do nothing
// }

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
