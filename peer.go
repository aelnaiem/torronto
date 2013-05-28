package torronto

import (
	"fmt"
)

type Peer struct {
	currentState State
	peers        Peers
	host         String
	port         int
}

func (peer Peer) Insert(filename String) int {
	// add the file to the local node and update status and filelist
	// get file info
	addLocalFile(filename, info)

	// divide the file by chunks and push it out
	// to peers
}

func (peer Peer) Query(status Status) int {
	// print out status of files
}

func (peer Peer) Join() int {
	makeFileMap()
	fileList = HostStatus.getFileList()

	//create the message to join
	joinMessage := encodeMessage(peer.host, peer.port, Join, fileList)
	// (add a timeout to update status to disconnected?)
	sendToAll(joinMessage)
}

func (peer Peer) Leave() int {
	// push out unique chunks, least replicated first

	// send out leave message
	leaveMessage := encodeMessage(peer.host, peer.port, Leave)
	sendToAll(leaveMessage)
}

func sendFileList(hostName string, portNumber int) {
	HostStatus.getFileList()
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
	sendMessage(hostName, portNumber, filesMessage)
}

func downloadFile(file File, conn net.Conn) {
	// check if we want to download the file and if we do:
	if f, ok := HostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunkNumber[0]] {
			// we already have the file, return?
		}
	}
	// save the content to the file by following
	// the path in the file name (file.Name:file.chunks[0])

	// update the status object

	// if we now have all the chunks, make the complete non-hidden file
}

func uploadFile(hostName string, portNumber int, file File) {
	// check if you have the file
	if f, ok := HostStatus.files[file.fileName]; ok {
		if f.chunks[file.chunkNumber[0]] {
			fileList := [1]File{file}
			filesMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

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
