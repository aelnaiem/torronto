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
	//error

	if action == Upload {
		tempMessage := make([]byte, HeaderSize, HeaderSize)
		copy(tempMessage, jsonMessage)
		jsonMessage = tempMessage
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
