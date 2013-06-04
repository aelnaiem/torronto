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
	NumFiles                 int
	Files                    []string
	Local                    []float32
	System                   []float32
	LeastReplication         []int
	WeightedLeastReplication []float32
}

func (status Status) numberofFiles() int {
	return len(status.replication)
}

func (status Status) fractionPresentLocally(fileArray []string) []float32 {
	fplArray := make([]float32, 0, len(fileArray))
	for file := range fileArray {
		local := float32(0.0)
		if f, ok := status.status["local"].files[fileArray[file]]; ok {
			for chunk := range f.Chunks {
				if f.Chunks[chunk] == 1 {
					local += 1
				}
			}
		}
		length := float32(len(status.status["local"].files[fileArray[file]].Chunks))
		fplArray = append(fplArray, local/length)
	}
	return fplArray
}

func (status Status) fractionPresent(fileArray []string) []float32 {
	fpArray := make([]float32, 0, len(fileArray))

	for file := range fileArray {
		levelZero := status.replication[fileArray[file]][0]
		missing := 0
		for chunk := range levelZero {
			if levelZero[chunk] == 1 {
				missing += 1
			}
		}
		fPresent := float32((len(levelZero) - missing)) / float32(len(levelZero))
		fpArray = append(fpArray, fPresent)
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
		length := float32(len(status.status["local"].files[fileArray[file]].Chunks))
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
	if _, ok := status.status[fullName].files[file.FileName]; !ok {
		chunks := make([]int, file.Chunks[0])
		for chunk := range chunks {
			chunks[chunk] = 0
		}
		chunks[file.Chunks[1]] = 1
		status.status[fullName].files[file.FileName] = File{
			FileName: file.FileName,
			Chunks:   chunks,
		}
	}
	status.status[fullName].files[file.FileName].Chunks[file.Chunks[1]] = 1

	incrementChunkReplication(file.FileName, file.Chunks[1], file.Chunks[0])
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
		status.status[fullName].files[file.FileName] = file
		for chunk := range file.Chunks {
			if file.Chunks[chunk] == 1 {
				incrementChunkReplication(file.FileName, chunk, len(file.Chunks))
				f := File{
					FileName: file.FileName,
					Chunks:   []int{len(file.Chunks), chunk},
				}
				localPeer.requestFile(f)
			}
		}
	}
}

func trackNewFile(file File) {
	status.status["local"].files[file.FileName] = file
	status.replication[file.FileName] = make([][]int, MaxPeers+1)
	status.replication[file.FileName][0] = file.Chunks
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
		for chunk := range file.Chunks {
			if file.Chunks[chunk] == 1 {
				decrementChunkReplication(file.FileName, chunk, len(file.Chunks))
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
