package torronto

import (
	"fmt"
	"math"
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
		var fileName string
		var chunkNumber int
		_, err := fmt.Sscanf(path, "%s:%d", &fileName, &chunkNumber)
		if err != nil {
			// error
		}

		if file, ok := hostStatus.files[fileName]; ok {
			file.chunks[chunkNumber] = 1
		} else {
			// return error
		}
	} else {
		numberOfChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))
		chunks := make([]int, numberOfChunks)
		for _, chunk := range chunks {
			chunk = 1
		}

		newFile := File{
			fileName: path,
			chunks:   chunks,
		}
		hostStatus.files[path] = newFile
		hostStatus.numFiles++
	}
	return nil
}

func makeFileList() {
	err := filepath.Walk("files", addLocalFile)
	if err != nil {
		// error
	}
}
