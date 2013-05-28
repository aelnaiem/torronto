package torronto

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"json"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>", os.Args[0])
		os.Exit(1)
	}

	addr := os.Args[1]
	addrArr := strings.Split(addr, ':')
	hostName, portNumber, err := addrAddr[0], Atoi(addrArr[1])
	if len(peerData) != 2 {
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
	go listenForMessages()

	// instantiate peer and join network
	peers := Peers{}
	peers.initialize()

	HostPeer = Peer{
		currentState: Connected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	// create status object
	HostStatus := Status{}

	// join the network
	peer.Join()
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
			HostPeer.Query(&status)
		}
	}
}

func listenForFiles() {
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				HostPeer.Insert(ev.Name)
			}
		case err := <-watcher.Error:
			// error
		}
	}
}

func listenForMessages() {
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
func sendMessage(hostName string, portNumber string, msg []byte) {
	// possibly want to add timeouts?
	ipAddresses, err := LookupIP(hostName)
	service := os.Args[1]

	service = net.TCPAddr{IP: ipAddresses[0], Port: port}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	_, err = conn.Write(msg)
	checkError(err)

	conn.Close()
}

func sendToAll(msg []byte) {
	for _, peer := range HostPeer.peers.peers {
		if !(peer.host == HostPeer.host && peer.port == HostPeer.port) {
			if peer.currentState != Disconnected {
				// add a timeout?
				sendMessage(p.host, p.port, leaveMessage)
			}
		}
	}
}

// handle a message from a peer
func handleMessage(conn net.Conn) {

	defer conn.Close()

	//read up to headerSize bytes
	jsonMessage = make([]byte, HeaderSize)
	for {
		n, err := conn.Read(message[0:])
		if err != nil {
			return
		}
	}

	// convert JSON message into type Message
	message := Message.decodeMessage(jsonMessage)

	// identify the type of message and act appropriately
	select {
	case message.action == Join:
		connectPeer(message.hostName, message.portNumber)
		sendFileList(message.hostName, message.portNumber)

	case message.action == Leave:
		disconnectPeer(message.hostName, message.portNumber)

	case message.action == Files:
		connectPeer(message.hostName, message.portNumber)
		updateStatus(message.files)

	case message.action == Upload:
		downloadFile(message.files[0], conn)

	case message.action == Download:
		uploadFile(message.hostName, message.portNumber, message.files[0])
	}
}

// handle the many different communication errors that can happen
func checkError(err error) {
	if err != nil {
		// error handling
	}
}
