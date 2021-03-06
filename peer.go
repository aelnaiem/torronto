package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Peer struct {
	currentState int
	peers        Peers
	host         string
	port         int
}

func (peer *Peer) insert(fileName string) {
	time.Sleep(500 * time.Millisecond)
	if _, ok := status.status["local"].files[fileName]; ok {
		return
	}
	reader, err := os.Open(fileName)
	for {
		if err == nil {
			reader.Close()
			break
		}
	}
	info, err := os.Stat(fileName)
	checkError(err)
	addLocalFile(fileName, info, nil)

	connectedPeers := make([]Peer, peer.peers.numPeers)
	for _, p := range peer.peers.peers {
		if p.currentState == Connected && (p.host != peer.host || p.port != peer.port) {
			connectedPeers = append(connectedPeers, p)
		}
	}

	if len(connectedPeers) == 0 {
		return
	}

	numChunks := int(math.Floor(float64(info.Size())/(ChunkSize+1) + 1))
	max := math.Max(float64(peer.peers.numPeers), float64(numChunks))
	chunk := 0
	p := 0
	for i := 0; i < int(max); {
		if chunk >= numChunks {
			chunk = 0
		}
		if p >= peer.peers.numPeers {
			p = 0
		}
		nextPeer := connectedPeers[p]
		peer.sendPeerChunk(nextPeer.host, nextPeer.port, fileName, numChunks, chunk, false)

		chunk += 1
		i += 1
		p += 1
	}

	for c := 0; c < numChunks; c++ {
		f := File{
			FileName: fileName,
			Chunks:   []int{numChunks, c},
		}
		fileList := []File{f}
		haveMessage := encodeMessage(peer.host, peer.port, Have, fileList)
		sendToAll(haveMessage)
	}

	return
}

func (peer Peer) query() {
	fileArray := make([]string, 0, len(status.replication))
	for file := range status.replication {
		fileArray = append(fileArray, file)
	}
	query := StatusInterface{
		NumFiles:                 status.numberofFiles(),
		Files:                    fileArray,
		Local:                    status.fractionPresentLocally(fileArray),
		System:                   status.fractionPresent(fileArray),
		LeastReplication:         status.minimumReplicationLevel(fileArray),
		WeightedLeastReplication: status.averageReplicationLevel(fileArray),
	}
	statusMessage, err := json.Marshal(query)
	checkError(err)
	fmt.Fprintf(os.Stderr, "Response: \n%s \n\n", statusMessage)
	return
}

func (peer *Peer) join() {
	peer.currentState = Connected
	fileList := status.getFileList()
	joinMessage := encodeMessage(peer.host, peer.port, Add, nil)

	jsonFiles, err := json.Marshal(fileList)
	checkError(err)

	tmp := make([]byte, ChunkSize, ChunkSize)
	copy(tmp, jsonFiles)
	jsonFiles = tmp

	messageToSend := append(joinMessage, jsonFiles...)
	sendToAll(messageToSend)
	return
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := status.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, nil)

	jsonFiles, err := json.Marshal(fileList)
	checkError(err)

	tmp := make([]byte, ChunkSize, ChunkSize)
	copy(tmp, jsonFiles)
	jsonFiles = tmp

	messageToSend := append(filesMessage, jsonFiles...)
	sendMessage(hostName, portNumber, messageToSend)
	return
}

func (peer *Peer) leave() {
	peer.currentState = Disconnected
	peer.peers.numPeers = 0
	for i := range peer.peers.peers {
		peer.peers.peers[i].currentState = Unknown
	}

	files := status.status["local"].files
	for file := range files {
		for chunk := range files[file].Chunks {
			if files[file].Chunks[chunk] == 1 {
				if status.replication[file][0][chunk] == 1 {
					peer.sendPeerChunk("", 0, file, len(status.replication[file][0]), chunk, true)
				}
			}
		}
	}

	leaveMessage := encodeMessage(peer.host, peer.port, Remove, nil)
	sendToAll(leaveMessage)
	peer.reset()
	return
}

func (peer *Peer) reset() {
	for peer := range status.status {
		if peer != "local" {
			delete(status.status, peer)
		}
	}
	// Clear replication data
	status.replication = make(map[string][][]int)

	// Update replication data for file chunks available locally
	files := status.status["local"].files
	for file := range files {
		for chunk := range files[file].Chunks {
			if files[file].Chunks[chunk] == 1 {
				incrementChunkReplication(file, chunk, len(files[file].Chunks))
			}
		}
	}
}

func (peer Peer) downloadFile(file File, conn *net.TCPConn) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			return
		}
	} else {
		chunks := make([]int, file.Chunks[0])
		for chunk := range chunks {
			chunks[chunk] = 0
		}
		chunks[file.Chunks[1]] = 1
		status.status["local"].files[file.FileName] = File{
			FileName: file.FileName,
			Chunks:   chunks,
		}
	}
	status.status["local"].files[file.FileName].Chunks[file.Chunks[1]] = 1
	incrementChunkReplication(file.FileName, file.Chunks[1], file.Chunks[0])

	err := conn.SetReadBuffer(ChunkSize)
	checkError(err)

	readBuffer := make([]byte, ChunkSize)
	_, err = conn.Read(readBuffer)
	checkError(err)
	conn.Close()

	status.mu.Lock()
	basepath := path.Dir(file.FileName)
	fileName := path.Base(file.FileName)
	err = os.MkdirAll(basepath, 0777)
	checkError(err)

	filePath := path.Join(basepath, fileName)

	localFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		for {
			localFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
			if err == nil {
				break
			}
		}
	}

	writeOffset := int64(file.Chunks[1] * ChunkSize)
	_, err = localFile.WriteAt(bytes.TrimRight(readBuffer, "\x00"), writeOffset)
	checkError(err)

	err = localFile.Close()
	checkError(err)

	status.mu.Unlock()
	fmt.Printf("Downloaded file %s:%d \n\n", file.FileName, file.Chunks[1])

	fileList := []File{file}
	haveMessage := encodeMessage(peer.host, peer.port, Have, fileList)
	sendToAll(haveMessage)
	return
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			peer.sendPeerChunk(hostName, portNumber, file.FileName, file.Chunks[0], file.Chunks[1], false)
			return
		}
	}

	fileList := []File{file}
	downloadMessage := encodeMessage(hostName, portNumber, Download, fileList)
	for _, p := range peer.peers.peers {
		if p.currentState == Connected && (p.host != peer.host || p.port != peer.port) {
			fullName := strings.Join([]string{p.host, strconv.Itoa(p.port)}, ":")
			if f, ok := status.status[fullName].files[file.FileName]; ok {
				if f.Chunks[file.Chunks[1]] == 1 {
					sendMessage(p.host, p.port, downloadMessage)
					return
				}
			}
		}
	}
	return
}

func (peer Peer) sendPeerChunk(hostName string, portNumber int, fileName string, numChunks int, chunk int, all bool) {
	f := File{
		FileName: fileName,
		Chunks:   []int{numChunks, chunk},
	}
	fileList := []File{f}
	uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

	writeBuffer := make([]byte, ChunkSize)
	readOffset := int64(chunk * ChunkSize)
	fileReading, err := os.Open(fileName)
	checkError(err)

	defer func() {
		if err := fileReading.Close(); err != nil {
			checkError(err)
		}
	}()

	_, err = fileReading.ReadAt(writeBuffer, readOffset)
	checkError(err)

	messageToSend := append(uploadMessage, writeBuffer...)
	if all {
		sendToAll(messageToSend)
	} else {
		sendMessage(hostName, portNumber, messageToSend)
	}
}

func (peer Peer) requestFile(hostName string, portNumber int, file File) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			return
		}
	}
	fileList := []File{file}
	fmt.Printf("Requesting file %s:%d\n\n", file.FileName, file.Chunks[1])
	downloadMessage := encodeMessage(peer.host, peer.port, Download, fileList)
	sendMessage(hostName, portNumber, downloadMessage)
	return
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
