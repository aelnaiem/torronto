package torronto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Peers struct {
	peers []Peer
}

func (peers Peers) initialize(peersFile string) {
	var hostName string
	var portNumber int

	content, err := ioutil.ReadFile("peersFile")
	checkError(err)

	lines := strings.Split(string(content), "\n")
	if len(lines) > MaxPeers {
		// TODO: either cut off extra list members or exit
	}

	peers.peers = make([]Peer, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		_, err := fmt.Sscanf(line, "%s %d", &hostName, &portNumber)
		checkError(err)

		peers.peers[i] = Peer{
			currentState: Unknown,
			host:         hostName,
			port:         portNumber,
		}
	}
	return
}

func (peers Peers) getPeer(hostName string, portNumber int) (Peer, error) {
	for _, peer := range peers.peers {
		if peer.host == hostName && peer.port == portNumber {
			return peer, nil
		}
	}
	return Peer{}, errors.New("Invalid host and/or port")
}

func (peers Peers) connectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)
	peer.currentState = Connected
}

func (peers Peers) disconnectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)
	peer.currentState = Disconnected
}