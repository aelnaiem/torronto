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

func (peer Peer) downloadFile(file File, conn *net.Conn) { //points to that connection
	// check if we want to download the file and if we do:
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			// we already have the file, TODO: close connection? return?
		}
	}
	//incorrect but we need to use TCPConn
	var tcpCon *net.TCPConn = TCPConn(conn)
	//unsure of the impact of setReadBuffer
	if err = tcpCon.SetReadBuffer(ChunkSize); err != nil { //set read buffer
		//error in setting read buffer size
	}
	readBuffer := make([]byte, ChunkSize)
	if readData, err := tcpCon.Read(readBuffer); err != nil {
		//error in reading from connection
	}

	dir, fileNameStr := Split(file.fileName)
	basepath := path.Dir(dir)
	filename := path.Base(fileNameStr)
	//creates the directory if it's not already there
	if os.MkdirAll(basepath, 0777) != nil {
		//error: panic("Unable to create directory for tagfile!")
	}

	// Create the tagfile.
	fileCreated, err := os.Create(path.Join([]string{basepath, filename}))
	//error check

	var writeOffset int64
	writeOffset = file.chunks[0] * ChunkSize
	if numberOfBytesWritten, err := fileCreated.WriteAt(readData, writeOffset); err != nil {
		//error check
	}

	//var walkFn WalkFunc{path:dir, info: nil, err: err }
	//walkFn = WalkFunc{}	//WalkFunc func(path string, info os.FileInfo, err error) error
	//set root
	//err := Walk(root, walkFn)
	//error check

	// update the status object
	hostStatus.files[file.fileName].chunks[file.chunks[0]] = 1
}

func (peer Peer) uploadFile(hostName string, portNumber int, file File) {
	// check if you have the file
	if f, ok := hostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunks[0]] == 1 {
			fileList := []File{file}
			uploadMessage := encodeMessage(peer.host, peer.port, Upload, fileList)
		}
		//incorrect but we need to use TCPConn
		var tcpCon *net.TCPConn = TCPConn(conn)
		//not sure of the impact of setWriteBuffer

		if err = tcpCon.SetWriteBuffer(ChunkSize); err != nil { //set write buffer
			//error in setting write buffer size
		}
		writeBuffer := make([]byte, ChunkSize)

		var readOffset int64
		readOffset = file.chunks[0] * ChunkSize

		var fileReading *os.File
		fileReading, err = Open(file.fileName)

		if numberOfBytesRead, err := fileReading.ReadAt(writeBuffer, readOffset); err != nil {
			//error check
		}

		if writeData, err := tcpCon.Write(writeBuffer); err != nil {
			//error in reading from connection
		}
	}
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
