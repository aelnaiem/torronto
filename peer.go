package main

import (
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

	file, err := os.Open(fileName)
	checkError(err)

	// TODO: why is this defered?
	defer func() {
		if err := file.Close(); err != nil {
			checkError(err)
		}
	}()

	// TODO: send each connected peer a file chunk

	// OLD CODE:
	// reader := bufio.NewReader(file)
	// buf := make([]byte, ChunkSize)
	// numberOfChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))

	// //TODO: iterate on hostName and port Number instead of chunks
	// for chunk := 0; chunk < numberOfChunks; chunk++ {
	// 	content, err := reader.Read(buf)
	// 	if err != io.EOF {
	// 		checkError(err)
	// 	}
	// 	if content == 0 {
	// 		break
	// 	}

	// 	f := File{
	// 		fileName: fileName,
	// 		chunks:   []int{chunk},
	// 	}
	// 	fileList := []File{f}
	// 	uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)
	// 	sendMessage(hostName, portNumber, uploadMessage, false)
	// }
}

func (peer Peer) query(hostName string, portNumber int) {
	status = Interface{
		"numFiles":                 status.numberofFiles(),
		"local":                    status.fractionPresentLocally(),
		"system":                   status.fractionPresent(),
		"leastReplication":         status.minimumReplicationLevel(),
		"weightedLeastReplication": status.averageReplicationLevel(),
	}

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
			fileList := []File{file}
			uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

			writeBuffer := make([]byte, ChunkSize)
			readOffset := int64(file.chunks[1] * ChunkSize)
			fileReading, err := os.Open(file.fileName)
			_, err = fileReading.ReadAt(writeBuffer, readOffset)
			checkError(err)

			messageToSend := append(uploadMessage, writeBuffer...)
			sendMessage(hostName, portNumber, messageToSend, false)
		}
	}
	return
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
