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
	// TODO: set status for all other peers for this file and update replication

	file, err := os.Open(filename)
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
	// 		fileName: filename,
	// 		chunks:   []int{chunk},
	// 	}
	// 	fileList := []File{f}
	// 	uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)
	// 	sendMessage(hostName, portNumber, uploadMessage, false)
	// }
}

func (peer Peer) Query(status *Status) {
	// TODO: print out status of files
	return
}

func (peer Peer) Join() {
	makeFileList()
	fileList := status.status["local"].getFileList()

	joinMessage := encodeMessage(peer.host, peer.port, Join, fileList)
	sendToAll(joinMessage, true)
	return
}

func (peer Peer) Leave() {
	// TODO: push out unique chunks, least replicated first

	leaveMessage := encodeMessage(peer.host, peer.port, Leave, nil)
	sendToAll(leaveMessage, false)
	return
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := status.status["local"].getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage, false)
	return
}

func (peer Peer) downloadFile(file File, tcpConn *net.TCPConn) {
	if f, ok := status.status["local"].files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			return
		}
	}

	_, err := tcpConn.SetReadBuffer(ChunkSize)
	checkError(err)

	readBuffer := make([]byte, ChunkSize)
	readBytes, err := tcpConn.Read(readBuffer)
	checkError(err)

	basepath := path.Dir(file.fileName)
	filename := path.Base(file.fileName)
	_, err = os.MkdirAll(basepath, 0777)
	checkError(err)

	filePath := path.Join(basepath, filename)
	localFile, err := os.OpenFile(filePath, os.O_CREAT|os.O_RDWR, 0777)
	checkError(err)

	writeOffset := int64(file.chunks[0] * ChunkSize)
	bytesWritten, err := localFile.WriteAt(readBuffer, writeOffset)
	checkError(err)

	status.status["local"].files[file.fileName].chunks[file.chunks[0]] = 1
	// TODO: send out Have message to all.
	return
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	if f, ok := status.status["local"].files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			fileList := []File{file}
			uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

			writeBuffer := make([]byte, ChunkSize)
			readOffset := int64(file.chunks[0] * ChunkSize)
			fileReading, err := os.Open(file.fileName) //var fileReading *os.File
			numberOfBytesRead, err := fileReading.ReadAt(writeBuffer, readOffset)
			checkError(err)

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
