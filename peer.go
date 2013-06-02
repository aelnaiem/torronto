package main

import (
	"encoding/json"
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

func (peer Peer) insert(fileName string) {
	info, err := os.Stat(fileName)
	checkError(err)

	addLocalFile(fileName, info, nil)

	numChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))
	chunk := 0
	p := 0
	max := math.Max(float64(peer.peers.numPeers), float64(numChunks))

	for i := 0; i < int(max); i++ {
		if chunk == numChunks {
			chunk = 0
		}
		if p == peer.peers.numPeers {
			p = 0
		}
		nextPeer := peer.peers.peers[p]
		if nextPeer.host == peer.host && nextPeer.port == peer.port {
			continue
		}
		peer.sendPeerChunk(nextPeer.host, nextPeer.port, fileName, chunk, false)
		chunk += 1
		p += 1
	}
}

func (peer Peer) query(hostName string, portNumber int) {
	fileArray := make([]string, 0, len(status.replication))
	for file := range status.replication {
		fileArray = append(fileArray, file)
	}
	query := StatusInterface{
		numFiles:                 status.numberofFiles(),
		files:                    fileArray,
		local:                    status.fractionPresentLocally(fileArray),
		system:                   status.fractionPresent(fileArray),
		leastReplication:         status.minimumReplicationLevel(fileArray),
		weightedLeastReplication: status.averageReplicationLevel(fileArray),
	}
	statusMessage, err := json.Marshal(query)
	checkError(err)
	sendMessage(hostName, portNumber, statusMessage, false)
	return
}

func (peer Peer) join() {
	makeFileList()
	fileList := status.getFileList()

	joinMessage := encodeMessage(peer.host, peer.port, Add, fileList)
	sendToAll(joinMessage, true)
	return
}

func (peer Peer) leave() {
	// TODO: push out unique chunks, least replicated first
	files := status.status["local"].files
	for file := range files {
		for chunk := range files[file].chunks {
			if files[file].chunks[chunk] == 1 {
				if status.replication[file][0][chunk] == 1 {
					peer.sendPeerChunk("", 0, file, chunk, true)
				}
			}
		}
	}

	leaveMessage := encodeMessage(peer.host, peer.port, Remove, nil)
	sendToAll(leaveMessage, false)
	return
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := status.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage, false)
	return
}

func (peer Peer) downloadFile(file File, tcpConn *net.TCPConn) {
	if f, ok := status.status["local"].files[file.fileName]; ok {
		if f.chunks[file.chunks[1]] == 1 {
			return
		}
	}

	err := tcpConn.SetReadBuffer(ChunkSize)
	checkError(err)

	readBuffer := make([]byte, ChunkSize)
	_, err = tcpConn.Read(readBuffer)
	checkError(err)

	basepath := path.Dir(file.fileName)
	fileName := path.Base(file.fileName)
	err = os.MkdirAll(basepath, 0777)
	checkError(err)

	filePath := path.Join(basepath, fileName)
	localFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	checkError(err)

	writeOffset := int64(file.chunks[1] * ChunkSize)
	_, err = localFile.WriteAt(readBuffer, writeOffset)
	checkError(err)

	status.status["local"].files[file.fileName].chunks[file.chunks[1]] = 1
	incrementChunkReplication(file.fileName, file.chunks[1], file.chunks[0])

	fileList := []File{file}
	haveMessage := encodeMessage(peer.host, peer.port, Have, fileList)
	sendToAll(haveMessage, false)

	return
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	if f, ok := status.status["local"].files[file.fileName]; ok {
		if f.chunks[file.chunks[1]] == 1 {
			peer.sendPeerChunk(hostName, portNumber, file.fileName, file.chunks[1], false)
		}
	}
	return
}

func (peer Peer) sendPeerChunk(hostName string, portNumber int, fileName string, chunk int, all bool) {
	f := File{
		fileName: fileName,
		chunks:   []int{chunk},
	}
	fileList := []File{f}
	uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

	writeBuffer := make([]byte, ChunkSize)
	readOffset := int64(chunk * ChunkSize)
	fileReading, err := os.Open(fileName)

	defer func() {
		if err := fileReading.Close(); err != nil {
			checkError(err)
		}
	}()

	_, err = fileReading.ReadAt(writeBuffer, readOffset)
	checkError(err)

	messageToSend := append(uploadMessage, writeBuffer...)
	if all {
		sendToAll(messageToSend, false)
	} else {
		sendMessage(hostName, portNumber, messageToSend, false)
	}
}

func (peer Peer) requestFile(file File) {
	if file, ok := status.status["local"].files[file.fileName]; ok {
		if file.chunks[file.chunks[1]] == 1 {
			return
		}
	}
	fileList := []File{file}
	downloadMessage := encodeMessage(peer.host, peer.port, Download, fileList)
	sendToAll(downloadMessage, false)
	return
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
