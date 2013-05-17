package p1 // think of a better name

import (
	"fmt"
)

type Peer struct {
	currentState State
	peers        Peers
}

func (peer Peer) insert(filename String) int {
	// code
}

func (peer Peer) query(status Status) int {
	// code
}

func (peer Peer) join() int {
	// codde
}

func (peer Peer) leave() int {
	// code
}

type State int

const (
	connected    = iota
	disconnected = iota
	unknown      = iota
)
