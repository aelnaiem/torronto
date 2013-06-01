package torronto

import (
	"strconv"
	"strings"
)

type Status struct {
	status      map[string]peerStatus
	replication map[string][]int
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

	// TODO: check if replication for that file exists and update it
	// send request to download if we don't have the file
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
		// TODO: update replication of the file
		// send request to download if we don't have the file
	}
}

func trackNewFile(file File) {
	status.status["local"].files[file.fileName] = file
	status.replication[file.fileName] = file.chunks
}
