package main

import (
	"math"
	"os"
	"path/filepath"
)

type File struct {
	FileName string
	Chunks   []int
}

type FileList struct {
	Files []File
}

func addLocalFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	numberOfChunks := int(math.Ceil(float64(info.Size()) / ChunkSize))
	chunks := make([]int, numberOfChunks)
	for chunk := range chunks {
		chunks[chunk] = 1
	}

	file := File{
		FileName: path,
		Chunks:   chunks,
	}

	trackNewFile(file)
	return nil
}

func makeFileList() {
	err := filepath.Walk("files", addLocalFile)
	checkError(err)
	return
}
