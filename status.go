package torronto

type Status struct {
	numFiles                 int
	local                    []float32
	system                   []float32
	leastReplication         []int
	weightedLeastReplication []float32
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
