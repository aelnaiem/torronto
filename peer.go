package torronto

import (
	"bufio"
	"io"
	"math"
	"net"
	"os"
)

type Peer struct {
	currentState int
	peers        Peers
	host         string
	port         int
}

func (peer Peer) Insert(filename string) int {
	// add the file to the local node and update status and filelist
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

	// iterate on hostName and port Number instead of chunks
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

}

func (peer Peer) Query(status *Status) int {
	// TODO: print out status of files
}

func (peer Peer) Join() int {
	makeFileList()
	fileList := hostStatus.getFileList()

	//create the message to join
	joinMessage := encodeMessage(peer.host, peer.port, Join, fileList)
	sendToAll(joinMessage, true)
}

func (peer Peer) Leave() int {
	// TODO: push out unique chunks, least replicated first

	// send out leave message
	leaveMessage := encodeMessage(peer.host, peer.port, Leave, nil)
	sendToAll(leaveMessage, false)
}

func (peer Peer) sendFileList(hostName string, portNumber int) {
	fileList := hostStatus.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage, false)
}

func (peer Peer) downloadFile(file File, conn net.Conn) {
	// check if we want to download the file and if we do:
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			// we already have the file, TODO: close connection? return?
		}
	}
	// TODO: save the content to the file by following
	// the path in the file name (file.Name:file.chunks[0])

	// update the status object
	hostStatus.files[file.fileName].chunks[file.chunks[0]] = 1

	// TODO: if we now have all the chunks, make the complete non-hidden file
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	// check if you have the file
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			fileList := []File{file}
			uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

			//TODO:
			// find the file and chunk in the directory, add it to the message and
			// send .(file.fileName:file.chunks[0]) or chunk file.fileName to the
			// appropriate part
		}
	}
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
