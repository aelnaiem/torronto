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
	numberOfChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))
	chunks := make([]int, numberOfChunks)
	for _, chunk := range chunks {
		chunk = 1
	}

	newFile := File{
		fileName: path,
		chunks:   chunks,
	}
	status.status["local"].files[path] = newFile
	return nil
}

func makeFileList() {
	err := filepath.Walk("files", addLocalFile)
	checkError(err)
}
