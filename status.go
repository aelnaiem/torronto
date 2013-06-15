package main

import (
	"strconv"
	"strings"
	"sync"
)

type Status struct {
	status      map[string]peerStatus
	replication map[string][][]int
	mu          sync.Mutex
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
		sum := 0
		if f, ok := status.replication[fileArray[file]]; ok {
			for r := range f {
				numReplicated := 0
				for chunk := 0; chunk < len(f[r]); chunk++ {
					if f[r][chunk] == 1 {
						numReplicated += 1
					}
				}
				sum += r * numReplicated
			}
		}
		length := float32(len(status.status["local"].files[fileArray[file]].Chunks))
		arlArray = append(arlArray, float32(sum)/length)
	}
	return arlArray
}

func (status Status) getFileList() FileList {
	fileList := FileList{}
	for _, file := range status.status["local"].files {
		fileList.Files = append(fileList.Files, file)
	}
	return fileList
}

func updateHaveStatus(hostName string, portNumber int, file File) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")
	status.mu.Lock()
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
		status.status[fullName].files[file.FileName] = File{
			FileName: file.FileName,
			Chunks:   chunks,
		}
	}

	if status.status[fullName].files[file.FileName].Chunks[file.Chunks[1]] == 0 {
		status.status[fullName].files[file.FileName].Chunks[file.Chunks[1]] = 1
		incrementChunkReplication(file.FileName, file.Chunks[1], file.Chunks[0])
	}
	status.mu.Unlock()

	localPeer.requestFile(file)
	return
}

func updateStatus(hostName string, portNumber int, files []File) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")
	if _, ok := status.status[fullName]; !ok {
		status.status[fullName] = peerStatus{
			files: make(map[string]File),
		}
	}

	for _, file := range files {
		for chunk := range file.Chunks {
			if file.Chunks[chunk] == 1 {
				status.mu.Lock()
				if _, ok := status.status[fullName].files[file.FileName]; ok {
					if status.status[fullName].files[file.FileName].Chunks[chunk] == 1 {
						status.mu.Unlock()
						continue
					}
				}
				incrementChunkReplication(file.FileName, chunk, len(file.Chunks))
				chunks := make([]int, 2)
				chunks[0], chunks[1] = len(file.Chunks), chunk
				f := File{
					FileName: file.FileName,
					Chunks:   chunks,
				}
				localPeer.requestFile(f)
				status.mu.Unlock()
			}
		}
		status.status[fullName].files[file.FileName] = file
	}
	return
}

func trackNewFile(file File) {
	chunksCopy := make([]int, len(file.Chunks))
	copy(chunksCopy, file.Chunks)
	status.status["local"].files[file.FileName] = File{
		FileName: file.FileName,
		Chunks:   chunksCopy,
	}
	status.replication[file.FileName] = make([][]int, MaxPeers+1)
	for i := 0; i <= MaxPeers; i++ {
		status.replication[file.FileName][i] = make([]int, len(file.Chunks))
		for j := 0; j < len(file.Chunks); j++ {
			if i == 1 {
				status.replication[file.FileName][i][j] = 1
			} else {
				status.replication[file.FileName][i][j] = 0
			}
		}
	}
}

func incrementChunkReplication(fileName string, chunkNumber int, numChunks int) {
	if _, ok := status.replication[fileName]; !ok {
		status.replication[fileName] = make([][]int, MaxPeers+1)
		for i := 0; i <= MaxPeers; i++ {
			status.replication[fileName][i] = make([]int, numChunks)
			for j := 0; j < numChunks; j++ {
				status.replication[fileName][i][j] = 0
			}
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
	return
}

func decrementPeerReplication(hostName string, portNumber int) {
	fullName := strings.Join([]string{hostName, strconv.Itoa(portNumber)}, ":")
	if _, ok := status.status[fullName]; !ok {
		return
	}

	for _, file := range status.status[fullName].files {
		for chunk := range file.Chunks {
			if file.Chunks[chunk] == 1 {
				decrementChunkReplication(file.FileName, chunk, len(file.Chunks))
			}
		}
	}

	status.status[fullName] = peerStatus{
		files: make(map[string]File),
	}
	return
}

func decrementChunkReplication(fileName string, chunkNumber int, numChunks int) {
	if _, ok := status.replication[fileName]; !ok {
		status.replication[fileName] = make([][]int, MaxPeers+1)
		for i := 0; i <= MaxPeers; i++ {
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
	return
}
