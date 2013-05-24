package torronto

type Message struct {
	HostName   string
	PortNumber string
	Action     string
	Files      []string
	File       string
}

// Message should have a function to create a message given certain
// values

// Messages should be created to a specific 'header' size so it's easy
// to tell what a message is vs. what a chunk is (for uploads)
