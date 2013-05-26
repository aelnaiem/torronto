package torronto

import (
	"fmt"
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

	// read for input from stdin and do the appropriate action
	// e.g. query for file, insert a file
	go listenForInput()

	// listen for messages?
	go listenForMessages()

	// instantiate peer and join network
	peers := Peers{}
	peers.initialize()

	hostPeer = Peer{
		currentState: Connected,
		peers:        peers,
		host:         hostName,
		port:         portNumber,
	}

	peer.Join()
}

// Listen for user input to stdin
func listenForInput() {
	for {
		var input string
		_, err := fmt.Scanln(&input)

		if err != nil {
			inputArr := strings.Split(input, " ")
			if strings.ToLower(inputArr[0]) == "query" {
				if len(inputArr) == 2 {
					hostPeer.Insert(inputArr[1])
				} else {
					fmt.Fprintf(os.Stderr, "Usage: insert <filename>")
				}
			}
			// TODO: accept other commands, like query...
		}
	}
}

func listenForMessages() {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// handle any new messages
		go handleMessage(conn)
	}
}

// handle a message from a peer
func handleMessage(conn net.Conn) {

	defer conn.Close()

	//read up to headerSize bytes
	receivedMessage = make([]byte, HeaderSize)
	for {
		n, err := conn.Read(receivedMessage[0:]) 
		if err != nil {
			return
		}
	}

	m := Message.decode_message(receivedMessage)	//convert JSON message into type Message

	//identify the type of message it is and perform the corresponding action
	if m.action == join {
		//get hostName and portNumber of the peer
		//then connectPear is called to update status of this peer
		connectPeer(m.HostName, m.PortNumber)
	}

	if m.action == leave {
		//get hostName and portNumber of the peer
		//then connectPear is called to update status of this peer
		disconnectPeer(m.HostName, m.PortNumber)
	}

	if m.action == files {
		// update status..
	}
	if m.action == upload {

	}
	if m.action == download {

	}
}

// Possibly use a check error function to handle the many different
// communication errors that can happen
func checkError(err error) {
	if err != nil {
		// error handling
	}
}
