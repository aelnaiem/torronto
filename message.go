package torronto

type Message struct {
	HostName   string
	PortNumber string
	Action     int
	Files      []File
}


func create_message(hostName string, portNumber string, action int, files []File) []byte {

	newinformation := Message{
		HostName: hostName,		//hostName = 32bits
		PortNumber: portNumber,	//portNumber = 16 bits
		Action: action,			//store the action as a string or an int? action = 3 to 4 bits?
		Files: files,
	}
	newMessage, err := json.Marshal(newinformation)

	return newMessage
}

//header size consists of hostName, portNumber and action

