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

	peers = make([]Peer, len(lines))
	numPeers = 0

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

		peers[i] = Peer{
			currentState: Unknown,
			host:         hostName,
			port:         portNumber,
		}
		numPeers++
	}
}

func (peers Peers) GetPeer(i int) Peer {
	return peers[i]
}

func (peers Peers) Visit(i int) {
	// not sure what this is for...
}

func (peers Peers) connectPeer(hostName string, portNumber int) {
	//for each entry in peers []Peer, find the one with this hostName
	//and portNumber and change status to Connected
	for _, peer := range peers {
		if peer.host == hostName && peer.port == portNumber {
			peer.currentState = Connected
			break
		}
	}
}

func (peers Peers) disconnectPeer(hostName string, portNumber int) {
	//for each entry in peers []Peer, find the one with this hostName
	//and portNumber and change status to Disconnected
	for _, peer := range peers {
		if peer.host == hostName && peer.port == portNumber {
			peer.currentState = Disconnected
			break
		}
	}
}
