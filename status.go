package torronto

type Status struct {
	numFiles int
	files    map[string]File
}

func (status Status) NumberofFiles() int {
	// code
}

func (status Status) FractionPresentLocally(fileNumber int) float32 {
	// code
}

func (status Status) FractionPresent(fileNumber int) float32 {
	// code
}

func (status Status) MinimumReplicationLevel(fileNumber int) int {
	// code
}

func (status Status) AverageReplicationLevel(fileNumber int) float32 {
	// code
}

func (status Status) getFileList() []File {
	fileList := []File
	for _, file := range status.files {
		append(fileList, file)
	}
	return fileList
}

func updateStatus([]File) {
	// update the status based on the files

	// send requests for the files we don't have (download request)
		f = File{
			filename: filename,
			chunks:   [1]int{chunk},
		}
}