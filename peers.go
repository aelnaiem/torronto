package main

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type Peers struct {
	numPeers int
	peers    []Peer
}

func (peers *Peers) initialize(peersFile string) {
	content, err := ioutil.ReadFile("peersFile")
	checkError(err)

	lines := strings.Split(string(content), "\n")
	if len(lines) > MaxPeers {
		// TODO: either cut off extra list members or exit
	}

	peers.peers = make([]Peer, len(lines))
	peers.numPeers = 0
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		peerData := strings.Split(string(line), " ")
		if len(peerData) != 2 {
			// incorrectly formed data, exit?
		}

		hostName := peerData[0]
		portNumber, err := strconv.Atoi(peerData[1])
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

func (peers *Peers) connectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)
	peer.currentState = Connected
	peers.numPeers += 1
}

func (peers *Peers) disconnectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)
	peer.currentState = Disconnected
	peers.numPeers -= 1
}
