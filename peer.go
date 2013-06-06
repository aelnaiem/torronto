package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"path"
)

type Peer struct {
	currentState int
	peers        Peers
	host         string
	port         int
}

func (peer *Peer) insert(fileName string) {
	if _, ok := status.status["local"].files[fileName]; ok {
		return
	}
	fmt.Println("ruhoh")
	reader, err := os.Open(fileName)
	fmt.Println(err)
	fmt.Println(fileName)
	for {
		if err == nil {
			reader.Close()
			break
		}
		fmt.Printf("ruhoh")
	}

	info, err := os.Stat(fileName)
	checkError(err)

	addLocalFile(fileName, info, nil)

	numChunks := int(math.Floor(float64(info.Size())/ChunkSize + 1))
	max := math.Max(float64(peer.peers.numPeers), float64(numChunks))

	chunk := 0
	p := 0
	for i := 0; i <= int(max)+1; {
		if chunk == numChunks {
			chunk = 0
		}
		if p == peer.peers.numPeers+1 {
			p = 0
		}
		nextPeer := peer.peers.peers[p]
		if nextPeer.host == peer.host && nextPeer.port == peer.port {
			p += 1
			i += 1
			continue
		}
		if nextPeer.currentState == Connected {
			peer.sendPeerChunk(nextPeer.host, nextPeer.port, fileName, numChunks, chunk, false)
			chunk += 1
		}
		p += 1
		i += 1
	}
	return
}

func (peer Peer) query(hostName string, portNumber int) {
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
	sendMessage(hostName, portNumber, statusMessage)
	return
}

func (peer *Peer) join() {
	peer.currentState = Connected
	fileList := status.getFileList()

	joinMessage := encodeMessage(peer.host, peer.port, Add, fileList)
	sendToAll(joinMessage)
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
	status.replication = make(map[string][][]int)
	files := status.status["local"].files
	for file := range files {
		for chunk := range files[file].Chunks {
			if files[file].Chunks[chunk] == 1 {
				incrementChunkReplication(file, chunk, len(files[file].Chunks))
			}
		}
	}
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := status.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage)
	return
}

func (peer Peer) downloadFile(file File, tcpConn *net.TCPConn) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			return
		}
	}

	err := tcpConn.SetReadBuffer(ChunkSize)
	checkError(err)

	readBuffer := make([]byte, ChunkSize)
	_, err = tcpConn.Read(readBuffer)
	checkError(err)

	basepath := path.Dir(file.FileName)
	fileName := path.Base(file.FileName)
	err = os.MkdirAll(basepath, 0777)
	checkError(err)

	filePath := path.Join(basepath, fileName)
	localFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	checkError(err)

	if _, ok := status.status["local"].files[file.FileName]; !ok {
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
	writeOffset := int64(file.Chunks[1] * ChunkSize)
	_, err = localFile.WriteAt(bytes.Trim(readBuffer, "\x00"), writeOffset)
	checkError(err)

	incrementChunkReplication(file.FileName, file.Chunks[1], file.Chunks[0])

	fileList := []File{file}
	haveMessage := encodeMessage(peer.host, peer.port, Have, fileList)
	sendToAll(haveMessage)

	return
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			peer.sendPeerChunk(hostName, portNumber, file.FileName, file.Chunks[0], file.Chunks[1], false)
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

func (peer Peer) requestFile(file File) {
	if f, ok := status.status["local"].files[file.FileName]; ok {
		if f.Chunks[file.Chunks[1]] == 1 {
			return
		}
	}
	fileList := []File{file}
	downloadMessage := encodeMessage(peer.host, peer.port, Download, fileList)
	sendToAll(downloadMessage)
	return
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
