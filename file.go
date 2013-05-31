package torronto

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	fileName string
	chunks   []int
}

func addLocalFile(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, ":") {
		fileData := strings.Split(string(path), ":")
		fileName, chunkNumber, err := fileData[0], Atoi(fileData[1])
		if err != nil {
			// error
		}

		if file, ok := HostStatus.files[fileName]; ok {
			file.chunks[chunkNumber] = 1
		} else {
			// return error
		}
	} else {
		numberOfChunks = int(math.Ceil(info.size / ChunkSize))
		chunks = make([]bool, numberOfChunks)
		for _, chunk := range chunks {
			chunk = 1
		}

		newFile = File{
			fileName: fileString,
			chunks:   chunks,
		}
		HostStatus.files[path] = newFile
		HostStatus.numFiles++
	}
}

func makeFileList() {
	err := filepath.Walk("files", addFile)
	if err != nil {
		// error
	}
}
