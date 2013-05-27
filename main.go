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

	hostPeer = Peer{
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
			hostPeer.Query(&status)
		}
	}
}

func listenForFiles() {
	for {
		select {
		case ev := <-watcher.Event:
			select {
			case ev.IsCreate():
				hostPeer.Insert(ev.Name)
			case ev.IsDelete():
				// not necessary...
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
		// handle any new messages
		go handleMessage(conn)
	}
}

//use this function to send message to a specified
//host name and port number
//note: HostName is for another peer and host is for the current peer
func sendMessage(host HostName, port portNumber, msg []byte) {
	//a lot of this is from that article Ahmed sent Anshika
	ipAddresses []IP, err := LookupIP(host)
	//if len(os.Args) != 2 {
	//	fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
	//	os.Exit(1)
	//}
	service := os.Args[1]

	service = net.TCPAddr{IP: ipAddresses[0], Port: port}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	//checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//checkError(err)
	_, err = conn.Write(msg)
	//checkError(err)
	result, err := ioutil.ReadAll(conn)
	//checkError(err)
	fmt.Println(string(result))
	//os.Exit(0)


	func DialTCP(net string, laddr, raddr *TCPAddr) (c *TCPConn, err os.Error)

	func (c *TCPConn) Write(msg) (int, error)
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

	m := Message.decode_message(receivedMessage) //convert JSON message into type Message

	//identify the type of message it is and perform the corresponding action
	if m.action == Join {
		//get hostName and portNumber of the peer
		//then connectPear is called to update status of this peer
		connectPeer(m.HostName, m.PortNumber)
	}

	if m.action == Leave {
		//get hostName and portNumber of the peer
		//then connectPear is called to update status of this peer
		disconnectPeer(m.HostName, m.PortNumber)
	}

	if m.action == Files {
		// update status..
	}
	if m.action == Upload {

	}
	if m.action == Download {

	}
}

// Possibly use a check error function to handle the many different
// communication errors that can happen
func checkError(err error) {
	if err != nil {
		// error handling
	}
}
