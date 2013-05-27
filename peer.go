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

	// divide the file by chunks and push it out
	// to peers
}

func (peer Peer) Query(status Status) int {
	// not sure what this is for...
}

func (peer Peer) Join() int {
	// get local fileList (maybe from status object?)
	var fileList []File = makeFileList()

	//create the message to join (add a timeout?)
	joinMessage := encodeMessage(peer.host, peer.port, Join, fileList)
	sendToAll(joinMessage)
}

func (peer Peer) Leave() int {
	// push out unique chunks, least replicated first

	// send out leave message
	leaveMessage := encodeMessage(peer.host, peer.port, Leave)
	sendToAll(leaveMessage)

	// close all sockets (only necessary if we keep connections open)
}

func sendFileList(hostName string, portNumber int) {
	// refer to either status object, or makeFileList() in file.go
	// to obtain fileList
	filesMessage := encodeMessage(peer.host, peer.port, Files, fileList)
}

func uploadFile(hostName string, portNumber int, file File) {
	fileList := [1]File{File}
	filesMessage := encodeMessage(peer.host, peer.port, Upload, fileList)

	// find the file and chunk in the directory, add it to the message and
	// send .(file.fileName:file.chunks[0]) or chunk file.fileName to the
	// appropriate part
}

func downloadFile(file File, conn net.Conn) {
	// check if we want to download the file and if we do:
	// save the content to the file .(file.Name:file.chunks[0])

	// if we have all the chunks, make the complete non-hidden file
	// update the status object
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
