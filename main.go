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
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>\n", os.Args[0])
		os.Exit(1)
	}

	if !strings.Contains(os.Args[1], ":") {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>\n", os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]

	addrArr := strings.Split(addr, ":")
	hostName := addrArr[0]
	portNumber, err := strconv.Atoi(addrArr[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: %s <host:port>\n", os.Args[0])
		os.Exit(1)
	}

	if _, err := os.Stat("files"); os.IsNotExist(err) {
		os.Mkdir("files", 0777)
	}
	checkError(err)

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

func sendMessage(hostName string, portNumber int, msg []byte) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(hostName, strconv.Itoa(portNumber)))
	checkError(err)

	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		localPeer.peers.disconnectPeer(hostName, portNumber)
	}

	_, err = conn.Write(msg)
	if err != nil {
		localPeer.peers.disconnectPeer(hostName, portNumber)
	}

	conn.Close()
	conn = nil
	return err
}

func sendToAll(msg []byte) {
	for _, peer := range localPeer.peers.peers {
		if !(peer.host == localPeer.host && peer.port == localPeer.port) {
			if peer.currentState != Disconnected {
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
	case message.Action == Join:
		var response []byte
		if localPeer.currentState == Connected {
			//fmt.Printf("Not joining: %s \n\n", message)
			response = encodeError(ErrConnected)
		} else {
			//fmt.Printf("Joining: %s \n\n", message)
			localPeer.join()
			response = encodeError(ErrOK)
		}
		sendMessage(message.HostName, message.PortNumber, response)
		return

	case message.Action == Leave:
		var response []byte
		if localPeer.currentState == Disconnected {
			//fmt.Printf("Not Leaving: %s \n\n", message)
			response = encodeError(ErrDisconnected)
		} else {
			//fmt.Printf("Leaving: %s \n\n", message)
			localPeer.leave()
			response = encodeError(ErrOK)
		}
		sendMessage(message.HostName, message.PortNumber, response)
		return

	case message.Action == Query:
		//fmt.Printf("query %s \n", message)
		if localPeer.currentState == Disconnected {
			response := encodeError(ErrDisconnected)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}
		localPeer.query(message.HostName, message.PortNumber)
		return

	case message.Action == Insert:
		//fmt.Printf("insert %s \n", message)
		if localPeer.currentState == Disconnected {
			response := encodeError(ErrDisconnected)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}

		src := message.Files[0].FileName
		dstArr := []string{"files", path.Base(message.Files[0].FileName)}
		dst := strings.Join(dstArr, "/")

		if _, ok := status.status["local"].files[dst]; ok {
			response := encodeError(ErrFileExists)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}

		sfile, err := os.Open(src)
		if err != nil {
			response := encodeError(ErrFileMissing)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}
		defer sfile.Close()

		dfile, err := os.Create(dst)
		if err != nil {
			response := encodeError(ErrBadPermission)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}
		defer dfile.Close()
		io.Copy(dfile, sfile)

		response := encodeError(ErrOK)
		sendMessage(message.HostName, message.PortNumber, response)
		return
	// peer messages
	// only act when peer is connected
	case localPeer.currentState == Connected:
		switch {
		case message.Action == Add:
			localPeer.peers.connectPeer(message.HostName, message.PortNumber, message.Files)
			localPeer.sendFileList(message.HostName, message.PortNumber)
			//fmt.Printf("Connected: %s \n\nPeer data: %d\n\n", message, localPeer.peers.peers)
			return

		case message.Action == Remove:
			localPeer.peers.disconnectPeer(message.HostName, message.PortNumber)
			//fmt.Printf("Disconnected: %s \n\nPeer data: %d\n\n", message, localPeer.peers.peers)
			return

		case message.Action == Files:
			localPeer.peers.connectPeer(message.HostName, message.PortNumber, message.Files)
			updateStatus(message.HostName, message.PortNumber, message.Files)
			//fmt.Printf("Connected: %s \n\nPeer data: %d\n\n", message, localPeer.peers.peers)
			return

		case message.Action == Upload:
			localPeer.downloadFile(message.Files[0], conn)
			return

		case message.Action == Download:
			localPeer.uploadFile(message.HostName, message.PortNumber, message.Files[0])
			return
		case message.Action == Have:
			updateHaveStatus(message.HostName, message.PortNumber, message.Files[0])
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println(err)
	}
}
