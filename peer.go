package torronto

import (
	"fmt"
)

type Peer struct {
	currentState State
	peers        Peers
}

func (peer Peer) Insert(filename String) int {
	// code
}

func (peer Peer) Query(status Status) int {
	// code
}

func (peer Peer) Join() int {
	// code
}

func (peer Peer) Leave() int {
	// code
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
