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

	// listen for status query
	go listenForQuery()

	// listen for files added to files folder
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// error
	}
	go listenForFiles()

	err = watcher.Watch("files")
	if err != nil {
		// error
	}

	// listen for messages?
	go listenForMessages(listener)

	// instantiate peer and join network
	peers := Peers{}
	peers.initialize("peerList")

	hostPeer = Peer{
		currentState: Connected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	// create status object
	hostStatus = Status{
		numFiles: 0,
		files:    make(map[string]File),
	}

	// join the network
	hostPeer.Join()
}

// Listen for user input to stdin
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
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleMessage(conn)
	}
}

// use this function to send message to a specified host and port
// might be better to connect to each peer only once, and keep track
// of open connections, rather than dialing every times?
func sendMessage(hostName string, portNumber int, msg []byte, hasTimeout bool) {
	addr := []string{hostName, strconv.Itoa(portNumber)}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strings.Join(addr, ":"))
	checkError(err)

	//saw these two lines of code online for adding a timeout
	//for read and write but not sure if it works completely
	var conn *net.TCPConn
	err = conn.SetDeadline(time.Now().Add(Timeout))
	conn, err = net.DialTCP("tcp", nil, tcpAddr)

	_, err = conn.Write(msg)
	checkError(err)

	//for now we close the connection after the attempt to send message
	conn.Close()
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

// handle a message from a peer
func handleMessage(conn net.Conn) {

	defer conn.Close()

	//read up to headerSize bytes
	jsonMessage := make([]byte, HeaderSize)
	for {
		n, err := conn.Read(jsonMessage[0:])
		checkError(err)
		print(n)
	}

	// convert JSON message into type Message
	message := decodeMessage(jsonMessage)

	// identify the type of message and act appropriately
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

// handle the many different communication errors that can happen
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
