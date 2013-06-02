package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var localPeer Peer
var status Status

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>", os.Args[0])
		os.Exit(1)
	}

	addr := os.Args[1]

	addrArr := strings.Split(addr, ":")
	hostName := addrArr[0]
	portNumber, err := strconv.Atoi(addrArr[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>", os.Args[0])
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	peers := Peers{}
	peers.initialize("peerList")
	localPeer = Peer{
		currentState: Disconnected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	status = Status{
		status:      make(map[string]peerStatus),
		replication: make(map[string][][]int),
	}
	status.status["local"] = peerStatus{
		files: make(map[string]File),
	}

	makeFileList()
	go listenForFiles()

	listenForMessages(listener)
}

func listenForFiles() {
	watcher, err := fsnotify.NewWatcher()
	checkError(err)

	err = watcher.Watch("files")
	checkError(err)

	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				localPeer.insert(ev.Name)
			}
		case err := <-watcher.Error:
			checkError(err)
		}
	}
}

func listenForMessages(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		checkError(err)
		go handleMessage(conn)
	}
}

func sendMessage(hostName string, portNumber int, msg []byte) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(hostName, strconv.Itoa(portNumber)))
	checkError(err)

	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	_, err = conn.Write(msg)
	checkError(err)

	conn.Close()
	conn = nil
}

func sendToAll(msg []byte, timeout bool) {
	for _, peer := range localPeer.peers.peers {
		if !(peer.host == localPeer.host && peer.port == localPeer.port) {
			if peer.currentState != Disconnected {
				if timeout {
					localPeer.peers.disconnectPeer(peer.host, peer.port)
				}
				sendMessage(peer.host, peer.port, msg)
			}
		}
	}
}

func handleMessage(conn *net.TCPConn) {
	defer conn.Close()

	jsonMessage := make([]byte, HeaderSize)
	_, err := conn.Read(jsonMessage)
	checkError(err)

	message := decodeMessage(jsonMessage)
	switch {
	// interface messages
	case message.action == Join:
		localPeer.join()
	case message.action == Leave:
		localPeer.leave()
	case message.action == Query:
		localPeer.query(message.hostName, message.portNumber)
	case message.action == Insert:
		src := message.files[0].fileName

		dstArr := []string{"files", path.Base(message.files[0].fileName)}
		dst := strings.Join(dstArr, "/")

		sfile, err := os.Open(src)
		checkError(err)
		defer sfile.Close()

		dfile, err := os.Create(dst)
		checkError(err)
		defer dfile.Close()
		io.Copy(dfile, sfile)

	// peer messages
	case localPeer.currentState == Connected:
		switch {
		case message.action == Add:
			localPeer.peers.connectPeer(message.hostName, message.portNumber)
			localPeer.sendFileList(message.hostName, message.portNumber)

		case message.action == Remove:
			localPeer.peers.disconnectPeer(message.hostName, message.portNumber)
			decrementPeerReplication(message.hostName, message.portNumber)

		case message.action == Files:
			localPeer.peers.connectPeer(message.hostName, message.portNumber)
			updateStatus(message.hostName, message.portNumber, message.files)

		case message.action == Upload:
			localPeer.downloadFile(message.files[0], conn)

		case message.action == Download:
			localPeer.uploadFile(message.hostName, message.portNumber, message.files[0])

		case message.action == Have:
			updateHaveStatus(message.hostName, message.portNumber, message.files[0])
		}
	}
}

func checkError(err error) {
	if err != nil {
		if err == io.EOF {

		} else if neterr, ok := err.(net.Error); ok && neterr.Timeout() {

		} else {

		}

	}
}
