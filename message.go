package torronto

type Message struct {
	HostName   string
	PortNumber string
	Action     string
	Files      []string
	File       string
}



func create_message(hostName string, portNumber string, action string, files []string, file string) []byte, int32 {
	//enter the information into the variable newMessage of type struct
	newinformation := Message{
		HostName: hostName,
		PortNumber: portNumber,
		Action: action,
		Files: files,
		File: file,
	}
	newMessage, err := json.Marshal(newinformation)

	// return the header size so that it's easy
	// to tell apart the header and the chunk/body
	//not sure if commas, spaces and colons need to be included
	headerLength = len("HostName") + len(hostName) + len("PortNumber") + len(portNumber) + len("Action") + len(action)

	return newMessage, headerLength
}

