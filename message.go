package torronto

type Message struct {
	hostName   string
	portNumber int
	action     int
	files      []File
}

func encodeMessage(hostName string, portNumber int, action int, files []File) []byte {
	message := Message{
		hostName:   hostName,   //hostName = 32bits
		portNumber: portNumber, //portNumber = 16 bits
		action:     action,     //store the action as a string or an int? action = 3 to 4 bits?
		files:      files,
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
