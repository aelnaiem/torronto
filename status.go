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

type StatusInterface struct {
	numFiles                 int
	files                    []string
	local                    []float32
	system                   []float32
	leastReplication         []int
	weightedLeastReplication []float32
}

func (status Status) numberofFiles() int {
	return len(status.replication)
}

func (status Status) fractionPresentLocally(fileArray []string) []float32 {
	fplArray := make([]float32, 0, len(fileArray))
	for file := range fileArray {
		local := float32(0.0)
		if f, ok := status.status["local"].files[fileArray[file]]; ok {
			for chunk := range f.chunks {
				if f.chunks[chunk] == 1 {
					local += 1
				}
			}
		}
		length := float32(len(status.status["local"].files[fileArray[file]].chunks))
		fplArray = append(fplArray, local/length)
	}
	return fplArray
}

func (status Status) fractionPresent(fileArray []string) []float32 {
	// TODO: figure out what the trick is here
	fpArray := make([]float32, 0, len(fileArray))
	for i := 0; i < len(fileArray); i++ {
		fpArray = append(fpArray, 1.0)
	}
	return fpArray
}

func (status Status) minimumReplicationLevel(fileArray []string) []int {
	mrlArray := make([]int, 0, len(fileArray))

	for file := range fileArray {
		lowest := 0
		if f, ok := status.replication[fileArray[file]]; ok {
		Search:
			for r := range f {
				for chunk := 0; chunk < len(f[r]); chunk++ {
					if f[r][chunk] == 1 {
						lowest = r
						break Search
					}
				}
			}
		}
		mrlArray = append(mrlArray, lowest)
	}
	return mrlArray
}

func (status Status) averageReplicationLevel(fileArray []string) []float32 {
	arlArray := make([]float32, 0, len(fileArray))

	for file := range fileArray {
		sum := float32(0)
		if f, ok := status.replication[fileArray[file]]; ok {
			for r := range f {
				numReplicated := 0
				for chunk := 0; chunk < len(f[r]); chunk++ {
					if f[r][chunk] == 1 {
						numReplicated += 1
					}
				}
				sum += float32(r * numReplicated)
			}
		}
		length := float32(len(status.status["local"].files[fileArray[file]].chunks))
		arlArray = append(arlArray, sum/length)
	}
	return arlArray
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
	status.replication[file.fileName] = make([][]int, MaxPeers+1)
	status.replication[file.fileName][0] = file.chunks
}

func incrementChunkReplication(fileName string, chunkNumber int, numChunks int) {
	if _, ok := status.replication[fileName]; !ok {
		status.replication[fileName] = make([][]int, MaxPeers+1)
		for i := 0; i < MaxPeers; i++ {
			status.replication[fileName][i] = make([]int, numChunks)
		}
	}

	replicationLevel := 0
	for i := MaxPeers; i >= 0; i-- {
		if status.replication[fileName][i][chunkNumber] == 1 {
			replicationLevel = i
			break
		}
	}

	status.replication[fileName][replicationLevel+1][chunkNumber] = 1
	status.replication[fileName][replicationLevel][chunkNumber] = 0
}

func decrementPeerReplication(hostName string, portNumber int) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")

	if _, ok := status.status[fullName]; !ok {
		status.status[fullName] = peerStatus{
			files: make(map[string]File),
		}
	}

	for _, file := range status.status[fullName].files {
		for chunk := range file.chunks {
			if file.chunks[chunk] == 1 {
				decrementChunkReplication(file.fileName, chunk, len(file.chunks))
			}
		}
	}
}

func decrementChunkReplication(fileName string, chunkNumber int, numChunks int) {
	if _, ok := status.replication[fileName]; !ok {
		status.replication[fileName] = make([][]int, MaxPeers+1)
		for i := 0; i < MaxPeers; i++ {
			status.replication[fileName][i] = make([]int, numChunks)
		}
	}

	replicationLevel := 0
	for i := 0; i <= MaxPeers; i++ {
		if status.replication[fileName][i][chunkNumber] == 1 {
			replicationLevel = i
			break
		}
	}

	status.replication[fileName][replicationLevel][chunkNumber] = 0
	status.replication[fileName][replicationLevel-1][chunkNumber] = 1
}
