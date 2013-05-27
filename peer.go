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
	filelist []File  := nil //shouldn't be nil, need to create the list of files here
	//create the message to join
	joinMessage:= encodeMessage(peer.host, peer.port, Join, fileList)


	// send out join message to all peers with your file list

	//why do this explicity?
	//if we get a response in time, set status to connected
	//peers.connectPeer(host, port)
	// else if they don't respond in time, set their status to
	//peers.disconnectPeer(host, port)
}

func (peer Peer) Leave() int {
	// push out unique chunks, least replicated first
	leaveMessage:= encodeMessage(peer.host, peer.port, Leave)
	// close all sockets
}

func sendFileList(hostName string, portNumber int) {
  // refer to either status object, or makeFileList() in file.go
  // to obtain fileList
  filesMessage:= encodeMessage(peer.host, peer.port, Files, fileList)
}

func sendFile(hostName string, portNumber int, file File) {
	fileList := [1]File{File}
	filesMessage:= encodeMessage(peer.host, peer.port, Upload, fileList)

	// find the file and chunk in the directory, add it to the message and
	// send .(file.fileName:file.chunks[0]) or chunk file.fileName to the
	// appropriate part
}

func downloadFile(file File, conn net.Conn) {

	// save the content to the file .(file.Name:file.chunks[0])

	// if we have all the chunks, make the complete non-hidden file
}

const (
	Connected    = iota
	Disconnected = iota
	Unknown      = iota
)
