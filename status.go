package p1

type Status struct {
	numFiles                 int
	local                    []float32
	system                   []float32
	leastReplication         []int
	weightedLeastReplication []float32
}

func (status Status) numberofFiles() int {
	// code
}

func (status Status) fractionPresentLocally(fileNumber int) float32 {
	// code
}

func (status Status) fractionPresent(fileNumber int) float32 {
	// code
}

func (status Status) minimumReplicationLevel(fileNumber int) int {
	// code
}

func (status Status) averageReplicationLevel(fileNumber int) float32 {
	// code
}
