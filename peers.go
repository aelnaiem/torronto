package torronto

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Peers struct {
	numPeers int
	peers    []Peer
}

func (peers Peers) Initialize(peersFile String) int {
	content, err := ioutil.ReadFile("peersFile")
	if err != nil {
		// exit gracefully if the file is not found
	}

	lines := strings.Split(string(content), "\n")
	if len(lines > MaxFiles) {
		// too many peers, should exit?
	}

	peers.peers = make([]Peer, len(lines))
	peers.numPeers = 0

	for i, line := range lines {
		// create a new peer and add it to the peers array and increment
		// the number of peers

		// empty line occurs at the end of the file
		if len(line) == 0 {
			continue
		}

		peerData := strings.Split(string(line), " ")
		if len(peerData) != 2 {
			// incorrectly formed data, exit?
		}

		hostName, portNumber, err := peerData[0], Atoi(peerData[1])
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
}

// GetPeer should probably take a hostName and portNumber
func (peers Peers) GetPeer(hostName string, portNumber int) (Peer, error) {
	for _, peer := range peers.peers {
		if peer.host == hostName && peer.port == portNumber {
			return peer, nil
		}
	}
	return nil, errors.New("invalid host and port")
}

func (peers Peers) ConnectPeer(hostName string, portNumber int) {
	peer, err = peers.GetPeer(hostName, portNumber)
	if err != nil {
		// couldn't find peer
	}
	peer.currentState = Connected
}

func (peers Peers) DisconnectPeer(hostName string, portNumber int) {
	peer, err = peers.GetPeer(hostName, portNumber)
	if err != nil {
		// couldn't find peer
	}
	peer.currentState = Disconnected
}
