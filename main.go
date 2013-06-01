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

var hostPeer Peer
var hostStatus Status

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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// error
	}
	go listenForFiles()

	err = watcher.Watch("files")
	if err != nil {
		// error
	}

	go listenForMessages(listener)

	peers := Peers{}
	peers.initialize("peerList")

	hostPeer = Peer{
		currentState: Connected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	hostStatus = Status{
		numFiles: 0,
		files:    make(map[string]File),
	}

	hostPeer.Join()
}

func listenForQuery() {
	for {
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			// error
		}

		inputArr := strings.Split(input, " ")
		if strings.ToLower(inputArr[0]) == "query" {
			hostPeer.Query(&hostStatus)
		}
	}
}

func listenForFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// error
	}
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				hostPeer.Insert(ev.Name)
			}
		case err := <-watcher.Error:
			print(err)
		}
	}
}

func listenForMessages(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go handleMessage(conn)
	}
}

func sendMessage(hostName string, portNumber int, msg []byte, hasTimeout bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(hostName, strconv.Itoa(portNumber)))
	checkError(err)

	var conn *net.TCPConn
	if hasTimeout == true {
		err = conn.SetDeadline(time.Now().Add(Timeout)) //does this work?
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)

	_, err = conn.Write(msg)
	checkError(err)

	conn.Close() //for now we close the connection after the attempt to send message
	conn = nil
}

func sendToAll(msg []byte, timeout bool) {
	for _, peer := range hostPeer.peers.peers {
		if !(peer.host == hostPeer.host && peer.port == hostPeer.port) {
			if peer.currentState != Disconnected {
				sendMessage(peer.host, peer.port, msg, timeout)
			}
		}
	}
}

func handleMessage(conn *net.TCPConn) {

	defer conn.Close()

	jsonMessage := make([]byte, HeaderSize)
	for {
		n, err := conn.Read(jsonMessage[0:])
		checkError(err)
		print(n)
	}

	message := decodeMessage(jsonMessage)
	switch {
	case message.action == Join:
		hostPeer.peers.connectPeer(message.hostName, message.portNumber)
		hostPeer.sendFileList(message.hostName, message.portNumber)

	case message.action == Leave:
		hostPeer.peers.disconnectPeer(message.hostName, message.portNumber)

	case message.action == Files:
		hostPeer.peers.connectPeer(message.hostName, message.portNumber)
		updateStatus(message.files)

	case message.action == Upload:
		hostPeer.downloadFile(message.files[0], conn)

	case message.action == Download:
		hostPeer.uploadFile(message.hostName, message.portNumber, message.files[0])
	}
}

func checkError(err error) {
	if err != nil {
		if err == io.EOF {
			//detected closed LAN connection
			//message not sent
		} else if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			//timout occurred
			//now what? do we give a timout error message or try
			//to resend it? or return a code to indicate that it should send
			//the message elsewhere?
		} else {
			// every other case
		}

	}
}
