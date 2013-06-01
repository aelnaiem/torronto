package torronto

type Status struct {
	numFiles int
	files    map[string]File
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
	for _, file := range status.files {
		fileList = append(fileList, file)
	}
	return fileList
}

func updateStatus([]File) {
	// TODO: update the status based on the files

	// TODO: send requests for the files we don't have (download request)
	/*f := File{
		fileName: filename,
		chunks:   []int{chunk},
	}*/
}
