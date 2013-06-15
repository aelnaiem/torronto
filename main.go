package main

import (
	"bufio"
	"code.google.com/p/go.exp/fsnotify"
	"fmt"
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
	go listenForCommand()

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

func listenForCommand() {
	for {
		in := bufio.NewReader(os.Stdin)
		s, err := in.ReadString('\n')
		checkError(err)
		input := strings.Split(s, " ")

		if strings.Contains(input[0], "join") {
			if localPeer.currentState == Connected {
				fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrConnected)
			} else {
				fmt.Printf("Joining... \n")
				localPeer.join()
				fmt.Printf("Joined \n")
			}
			fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrOK)
		} else if strings.Contains(input[0], "leave") {
			fmt.Printf("input: %s", input[0])
			if localPeer.currentState == Disconnected {
				fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrDisconnected)
			} else {
				fmt.Printf("Leaving... \n")
				localPeer.leave()
				fmt.Printf("Left \n")
			}
			fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrOK)
		} else if strings.Contains(input[0], "query") {
			if localPeer.currentState == Disconnected {
				fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrDisconnected)
			}
			localPeer.query()
		} else if strings.Contains(input[0], "insert") {
			if localPeer.currentState == Disconnected {
				fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrDisconnected)
			} else {
				src := strings.Trim(input[1], "\n")
				dstArr := []string{"files", path.Base(src)}
				dst := strings.Join(dstArr, "/")

				fmt.Printf("Inserting File: %s into %s... \n", src, dst)
				if _, ok := status.status["local"].files[dst]; ok {
					fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrFileExists)
				} else {

					sfile, err := os.Open(src)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrFileMissing)
					} else {
						defer sfile.Close()

						dfile, err := os.Create(dst)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrBadPermission)
						} else {
							defer dfile.Close()
							io.Copy(dfile, sfile)

							fmt.Fprintf(os.Stderr, "Response: %d \n\n", ErrOK)
							fmt.Printf("Inserted \n\n")
						}
					}
				}
			}
		}
	}
}

func listenForMessages(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err == nil && localPeer.currentState == Connected {
			go handleMessage(conn)
		} else {
			conn.Close()
		}
	}
}

func sendMessage(hostName string, portNumber int, msg []byte) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort(hostName, strconv.Itoa(portNumber)))
	checkError(err)

	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	_, err = conn.Write(msg)

	return
}

func sendToAll(msg []byte) {
	for _, peer := range localPeer.peers.peers {
		if !(peer.host == localPeer.host && peer.port == localPeer.port) {
			if peer.currentState != Disconnected {
				sendMessage(peer.host, peer.port, msg)
			}
		}
	}
	return
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
			fmt.Printf("Not joining: %s \n\n", message)
			response = encodeError(ErrConnected)
		} else {
			fmt.Printf("Joining: %s \n\n", message)
			localPeer.join()
			response = encodeError(ErrOK)
		}
		sendMessage(message.HostName, message.PortNumber, response)
		return

	case message.Action == Leave:
		var response []byte
		if localPeer.currentState == Disconnected {
			response = encodeError(ErrDisconnected)
		} else {
			localPeer.leave()
			response = encodeError(ErrOK)
		}
		sendMessage(message.HostName, message.PortNumber, response)
		return

	case message.Action == Query:
		if localPeer.currentState == Disconnected {
			response := encodeError(ErrDisconnected)
			sendMessage(message.HostName, message.PortNumber, response)
			return
		}
		localPeer.query()
		return

	case message.Action == Insert:
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
	case message.Action == Add:
		localPeer.peers.connectPeer(message.HostName, message.PortNumber, conn)
		localPeer.sendFileList(message.HostName, message.PortNumber)
		fmt.Printf("Connected: %s:%d\n\n", message.HostName, message.PortNumber)
		return

	case message.Action == Remove:
		localPeer.peers.disconnectPeer(message.HostName, message.PortNumber)
		fmt.Printf("Disconnected: %s:%d\n\n", message.HostName, message.PortNumber)
		return

	case message.Action == Files:
		localPeer.peers.connectPeer(message.HostName, message.PortNumber, conn)
		fmt.Printf("Updated file list from: %s:%d\n\n", message.HostName, message.PortNumber)
		return

	case message.Action == Upload:
		localPeer.downloadFile(message.Files[0], conn)
		return

	case message.Action == Download:
		localPeer.uploadFile(message.HostName, message.PortNumber, message.Files[0])
		return
	case message.Action == Have:
		updateHaveStatus(message.HostName, message.PortNumber, message.Files[0])
		fmt.Printf("Updated status that %s:%d Has file:chunk %s:%d\n\n", message.HostName, message.PortNumber, message.Files[0].FileName, message.Files[0].Chunks[1])
		return
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
