package torronto

import (
	"encoding/json"
)

type Message struct {
	hostName   string
	portNumber int
	action     int
	files      []File
}

func encodeMessage(hostName string, portNumber int, action int, files []File) []byte {
	message := Message{
		hostName:   hostName,
		portNumber: portNumber,
		action:     action,
		files:      files,
	}
	jsonMessage, err := json.Marshal(message)
	checkError(err)

	if action == Upload {
		tmp := make([]byte, HeaderSize, HeaderSize)
		copy(tmp, jsonMessage)
		jsonMessage = tmp
	}
	return jsonMessage
}

func decodeMessage(jsonMessage []byte) Message {
	var message Message
	err := json.Unmarshal(jsonMessage, &message)
	checkError(err)

	return message
}

const (
	Join     = iota
	Leave    = iota
	Have     = iota
	Files    = iota
	Upload   = iota
	Download = iota
)
