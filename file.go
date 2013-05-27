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

func addFile(path string, info os.FileInfo, err error) error {
	// files object doesn't exist. Possibly tie it to the status
	// object...
	segments := strings.Split(path, string(filepath.Separator))
	fileString = segments[len(segments-1)]

	if strings.Contains(fileString, ":") {
		fileData := strings.Split(string(fileString), ":")
		fileName, chunkNumber, err := peerData[0], Atoi(peerData[1])
		if err != nil {
			// error
		}

		for i, file := range files {
			if file.fileName == fileName {
				append(file.chunks, chunkNumber)
				break
			}
			if i == files.length {
				chunks := [...]int{chunkNumber}
				newFile = File{
					fileName: fileName,
					chunks:   chunks,
				}
				append(files, newFile)
			}
		}
	} else {
		numberOfChunks = int(math.Ceil(info.size / ChunkSize))
		chunks = make([]int, numberOfChunks)
		for i, chunk := range chunks {
			chunk = i
		}

		newFile = File{
			fileName: fileString,
			chunks:   chunks,
		}
		files = append(files, newFile)
	}
}

// this should definitely be only run once
// after initial file list is set, it should be updated
// with new file data but not recreated from scratch with this function
// might just put this with status...
func makeFileList() []Files {
	err := filepath.Walk("files", addFile)
	if err != nil {
		// error
	}
}
