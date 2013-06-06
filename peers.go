package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Fprintf(os.Stderr, "Too many peers in peersFile\n")
		os.Exit(1)
	}

	peers.peers = make([]Peer, len(lines))
	peers.numPeers = 0
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		peerData := strings.Split(string(line), " ")
		if len(peerData) != 2 {
			continue
		}

		hostName := peerData[0]
		portNumber, err := strconv.Atoi(peerData[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Malformed peer listing in peersFile line: %d\n", i)
			os.Exit(1)
		}

		peers.peers[i] = Peer{
			currentState: Unknown,
			host:         hostName,
			port:         portNumber,
		}
	}
	return
}

func (peers Peers) getPeer(hostName string, portNumber int) (*Peer, error) {
	for i, peer := range peers.peers {
		if peer.host == hostName && peer.port == portNumber {
			return &peers.peers[i], nil
		}
	}
	return &Peer{}, errors.New("Invalid host and/or port")
}

func (peers *Peers) connectPeer(hostName string, portNumber int, files []File) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)

	if peer.currentState != Connected {
		updateStatus(hostName, portNumber, files)
		peers.numPeers += 1
		// fmt.Printf("Number of peers: %d (inc)\n\n", peers.numPeers)
	}

	peer.currentState = Connected
	return
}

func (peers *Peers) disconnectPeer(hostName string, portNumber int) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)

	if peer.currentState == Connected {

		decrementPeerReplication(hostName, portNumber)
		peers.numPeers -= 1
		fmt.Printf("Removing: %d\n\n", portNumber)
	}

	peer.currentState = Disconnected
	return
}
