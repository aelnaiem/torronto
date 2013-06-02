package main

import (
	"strconv"
	"strings"
)

type Status struct {
	status      map[string]peerStatus
	replication map[string][][]int
}

type peerStatus struct {
	files map[string]File
}

func (status Status) NumberofFiles() int {
	// TODO: code
	return 0
}

func (status Status) FractionPresentLocally(fileNumber int) float32 {
	// TODO: code
	return 0
}

func (status Status) FractionPresent(fileNumber int) float32 {
	// TODO: code
	return 0
}

func (status Status) MinimumReplicationLevel(fileNumber int) int {
	// TODO: code
	return 0
}

func (status Status) AverageReplicationLevel(fileNumber int) float32 {
	// TODO: code
	return 0
}

func (status Status) getFileList() []File {
	fileList := []File{}
	for _, file := range status.status["local"].files {
		fileList = append(fileList, file)
	}
	return fileList
}

func updateHaveStatus(hostName string, portNumber int, file File) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")

	if _, ok := status.status[fullName]; !ok {
		status.status[fullName] = peerStatus{
			files: make(map[string]File),
		}
	}
	if _, ok := status.status[fullName].files[file.fileName]; !ok {
		chunks := make([]int, file.chunks[0])
		for chunk := range chunks {
			chunks[chunk] = 0
		}
		chunks[file.chunks[1]] = 1
		status.status[fullName].files[file.fileName] = File{
			fileName: file.fileName,
			chunks:   chunks,
		}
	}
	status.status[fullName].files[file.fileName].chunks[file.chunks[1]] = 1

	incrementChunkReplication(file.fileName, file.chunks[1], file.chunks[0])
	localPeer.requestFile(file)
}

func updateStatus(hostName string, portNumber int, files []File) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")

	if _, ok := status.status[fullName]; !ok {
		status.status[fullName] = peerStatus{
			files: make(map[string]File),
		}
	}
	for _, file := range files {
		status.status[fullName].files[file.fileName] = file
		for chunk := range file.chunks {
			if file.chunks[chunk] == 1 {
				incrementChunkReplication(file.fileName, chunk, len(file.chunks))
				f := File{
					fileName: file.fileName,
					chunks:   []int{len(file.chunks), chunk},
				}
				localPeer.requestFile(f)
			}
		}
	}
}

func trackNewFile(file File) {
	status.status["local"].files[file.fileName] = file
	status.replication[file.fileName] = make([][]int, MaxPeers)
	status.replication[file.fileName][0] = file.chunks
}

func incrementChunkReplication(fileName string, chunkNumber int, numChunks int) {
	if _, ok := status.replication[fileName]; !ok {
		status.replication[fileName] = make([][]int, MaxPeers)
		for i := 0; i < MaxPeers; i++ {
			status.replication[fileName][i] = make([]int, numChunks)
		}
	}

	replicationLevel := 0
	for i := 0; i < MaxPeers; i++ {
		if status.replication[fileName][i][chunkNumber] == 1 {
			replicationLevel = i
			break
		}
	}

	status.replication[fileName][replicationLevel][chunkNumber] = 0
	status.replication[fileName][replicationLevel][chunkNumber] = 1
}
