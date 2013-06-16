package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
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
	if (len(lines) - 1) > MaxPeers {
		fmt.Fprintf(os.Stderr, "Too many peers in peersFile\n")
		os.Exit(1)
	}

	peersInFile := len(lines) - 1
	peers.peers = make([]Peer, peersInFile)
	peers.numPeers = 0
	for i, line := range lines {
		if i < peersInFile && i <= MaxPeers {
			if len(line) == 0 {
				continue
			}
			if strings.Contains(string(line), "\r") == true {
				line = strings.TrimRight(string(line), "\r")
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

func (peers *Peers) connectPeer(hostName string, portNumber int, conn *net.TCPConn) {
	peer, err := peers.getPeer(hostName, portNumber)
	checkError(err)

	err = conn.SetReadBuffer(ChunkSize)
	checkError(err)

	readBuffer := make([]byte, ChunkSize)
	_, err = conn.Read(readBuffer)
	checkError(err)

	var fileList FileList
	readBuffer = bytes.TrimRight(readBuffer, "\x00")

	err = json.Unmarshal(readBuffer, &fileList)
	checkError(err)

	if peer.currentState != Connected {
		updateStatus(hostName, portNumber, fileList.Files)
		peers.numPeers += 1
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
		fmt.Printf("Decrementing: %s:%d\n\n", hostName, portNumber)
	}

	peer.currentState = Disconnected
	return
}
