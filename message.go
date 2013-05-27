package torronto

type Message struct {
	HostName   string
	PortNumber int
	Action     int
	Files      []File
}

func encodeMessage(hostName string, portNumber int, action int, files []File) []byte {
	message := Message{
		HostName:   hostName,   //hostName = 32bits
		PortNumber: portNumber, //portNumber = 16 bits
		Action:     action,     //store the action as a string or an int? action = 3 to 4 bits?
		Files:      files,
	}
	jsonMessage, err := json.Marshal(message)

	//add padding to the message
	if len(jsonMessage) < HeaderSize {

	}

	return jsonMessage
}

func decodeMessage(jsonMessage []byte) Message {
	var message Message
	err := json.Unmarshal(jsonMessage, &message)

	return message
}

const (
	Join     = iota
	Leave    = iota
	Files    = iota
	Upload   = iota
	Download = iota
)
