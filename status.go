package torronto

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
	// TODO: check if peer exists in status object
	// TODO: check if file exists in peers files and add it

	// TODO: check if replication for that file exists and update it
}

func updateStatus(hostName string, portNumber int, files []File) {
	fullName := join([]string{hostName, strconv.Itoa(portNumber)}, ":")

	// TODO: check if peer exists in status object
	for _, file := range files {
		// TODO: check if file exists in the files
		if f, ok := status.status[fullName].files[file.fileName]; ok {
			status.status[fullName].files[f.fileName].chunks = f.chunks
			// TODO: check if replication for that file exists
			for chunk := range f.chunks {
				status.replication[f.fileName].chunks[chunk] += 1
			}
		}
		status.status[full]
	}
}

func trackNewFile(file File) {
	status.status["local"].files[path] = file
	status.replication[file.fileName] = file.chunks
}
