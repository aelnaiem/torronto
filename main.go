package torronto

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var localPeer Peer
var status Status

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>", os.Args[0])
		os.Exit(1)
	}

	var hostName string
	var portNumber int
	addr := os.Args[1]
	_, err := fmt.Sscanf(addr, "%s %d", &hostName, &portNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>", os.Args[0])
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	go listenForQuery()

	go listenForFiles()

	go listenForMessages(listener)

	peers := Peers{}
	peers.initialize("peerList")

	localPeer = Peer{
		currentState: Connected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	status = Status{
		status:      make(map[string]peerStatus),
		replication: make(map[string][]int),
	}
	status.status["local"] = peerStatus{
		files: make(map[string]File),
	}

	localPeer.Join()
}

func listenForQuery() {
	for {
		var input string
		_, err := fmt.Scanln(&input)
		checkError(err)

		inputArr := strings.Split(input, " ")
		if strings.ToLower(inputArr[0]) == "query" {
			localPeer.Query(&status.status["local"])
		}
	}
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
				localPeer.Insert(ev.Name)
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

func sendMessage(hostName string, portNumber int, msg []byte, hasTimeout bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(hostName, strconv.Itoa(portNumber)))
	checkError(err)

	var conn *net.TCPConn
	if hasTimeout == true {
		err = conn.SetDeadline(time.Now().Add(Timeout))
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)

	_, err = conn.Write(msg)
	checkError(err)

	conn.Close()
	conn = nil
}

func sendToAll(msg []byte, timeout bool) {
	for _, peer := range localPeer.peers.peers {
		if !(peer.host == localPeer.host && peer.port == localPeer.port) {
			if peer.currentState != Disconnected {
				sendMessage(peer.host, peer.port, msg, timeout)
			}
		}
	}
}

func handleMessage(conn *net.TCPConn) {

	defer conn.Close()

	jsonMessage := make([]byte, HeaderSize)
	n, err := conn.Read(jsonMessage[0:])
	checkError(err)

	message := decodeMessage(jsonMessage)
	switch {
	case message.action == Join:
		localPeer.peers.connectPeer(message.hostName, message.portNumber)
		localPeer.sendFileList(message.hostName, message.portNumber)

	case message.action == Leave:
		localPeer.peers.disconnectPeer(message.hostName, message.portNumber)

	case message.action == Files:
		localPeer.peers.connectPeer(message.hostName, message.portNumber)
		updateStatus(message.files)

	case message.action == Upload:
		localPeer.downloadFile(message.files[0], conn)

	case message.action == Download:
		localPeer.uploadFile(message.hostName, message.portNumber, message.files[0])
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
