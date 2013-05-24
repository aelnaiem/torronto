package torronto

import (
	"fmt"
	"io/ioutil"
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

  hostPeer = Peer {
		currentState: Connected,
		peers,				peers,
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
			if strings.ToLower(inputArr[0]) == "insert" {
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
 // TODO: translate message and then...
	if (new peer or filelist) {
	  // TODO: get hostName and portNumber and then...
	  hostPeer.Peers.connectPeer(hostName, portNumber)
	}

	if (leaving peer) {
	  hostPeer.Peers.disconnectPeer(hostName, portNumber)
	}

	if (filelist) {
		// update status..
	}
}

// Possibly use a check error function to handle the many different
// communication errors that can happen
func checkError(err error) {
	if err != nil {
		// error handling
	}
}
