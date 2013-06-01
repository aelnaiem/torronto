package torronto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Peers struct {
	numPeers int
	peers    []Peer
}

func (peers Peers) initialize(peersFile string) {
	var hostName string
	var portNumber int

	content, err := ioutil.ReadFile("peersFile")
	if err != nil {
		// exit gracefully if the file is not found
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > MaxPeers {
		// too many peers, should exit?
	}

	peers.peers = make([]Peer, len(lines))
	peers.numPeers = 0

	for i, line := range lines {
		// create a new peer and add it to the peers array and increment
		// the number of peers

		if len(line) == 0 {
			continue
		}

		_, err := fmt.Sscanf(line, "%s %d", &hostName, &portNumber)
		if err != nil {
			// portNumber given was not an integer, exit?
		}

		peers.peers[i] = Peer{
			currentState: Unknown,
			host:         hostName,
			port:         portNumber,
		}
		peers.numPeers++
	}
	return
}

func (peers Peers) getPeer(hostName string, portNumber int) (Peer, error) {
	for _, peer := range peers.peers {
		if peer.host == hostName && peer.port == portNumber {
			return peer, nil
		}
	}
	return Peer{}, errors.New("invalid host and port")
}

func (peers Peers) connectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	if err != nil {
		// couldn't find peer
	}
	peer.currentState = Connected
}

func (peers Peers) disconnectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	if err != nil {
		// couldn't find peer
	}
	peer.currentState = Disconnected
}
