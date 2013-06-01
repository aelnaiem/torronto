package torronto

import (
	"bufio"
	"io"
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

func (peer Peer) Insert(filename string) {
	info, err := os.Stat(filename)
	checkError(err)

	addLocalFile(filename, info, nil)

	file, err := os.Open(filename)

	if err != nil {
		// err
	}

	defer func() {
		if err := file.Close(); err != nil {
			// err
		}
	}()

	read := bufio.NewReader(file)
	buf := make([]byte, ChunkSize)
	numberOfChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))

	//TODO: iterate on hostName and port Number instead of chunks
	for chunk := 0; chunk < numberOfChunks; chunk++ {
		content, err := read.Read(buf)
		if err != io.EOF {
			checkError(err)
		}
		if content == 0 {
			break
		}
		f := File{
			fileName: filename,
			chunks:   []int{chunk},
		}
		fileList := []File{f}
		uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

		// TODO: append content to uploadmessage

		// sendMessage(hostName, portNumber, uploadMessage, false)
		// divide the file by chunks and push it out
		// to peers
	}
	return
}

func (peer Peer) Query(status *Status) {
	// TODO: print out status of files
	return
}

func (peer Peer) Join() {
	makeFileList()
	fileList := hostStatus.getFileList()

	joinMessage := encodeMessage(peer.host, peer.port, Join, fileList)
	sendToAll(joinMessage, true)
	return
}

func (peer Peer) Leave() {
	// TODO: push out unique chunks, least replicated first

	// send out leave message
	leaveMessage := encodeMessage(peer.host, peer.port, Leave, nil)
	sendToAll(leaveMessage, false)
	return
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := hostStatus.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage, false)
	return
}

func (peer Peer) downloadFile(file File, tcpConn *net.TCPConn) { //points to that connection
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			return
		}
	}

	if err := tcpConn.SetReadBuffer(ChunkSize); err != nil { //unsure of the impact of setReadBuffer
		//error
	}
	readBuffer := make([]byte, ChunkSize)
	readBytes, err := tcpConn.Read(readBuffer)
	if err != nil {
		//error in reading from connection
	}

	basepath := path.Dir(file.fileName)
	filename := path.Base(file.fileName)
	if os.MkdirAll(basepath, 0777) != nil {
		//error
	}

	fileCreated, err := os.Create(path.Join(basepath, filename))
	//error

	var writeOffset int64
	writeOffset = int64(file.chunks[0] * ChunkSize)
	if numberOfBytesWritten, err := fileCreated.WriteAt(readBuffer, writeOffset); err != nil {
		//error
	}

	hostStatus.files[file.fileName].chunks[file.chunks[0]] = 1
	return
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			fileList := []File{file}
			uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

			writeBuffer := make([]byte, ChunkSize)
			readOffset := int64(file.chunks[0] * ChunkSize)
			fileReading, err := os.Open(file.fileName) //var fileReading *os.File
			if numberOfBytesRead, err := fileReading.ReadAt(writeBuffer, readOffset); err != nil {
				//error
			}

			messageToSend := append(uploadMessage, writeBuffer...)
			sendMessage(hostName, portNumber, messageToSend, false)
		}
	}
	return
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
