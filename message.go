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

	//adds padding to message to make sure length = headerSize
	if action == Upload {
		//make a new array tempMessage with the desired length
		tempMessage := make([]byte, HeaderSize, HeaderSize)
		//copy contents of jsonMessage into tempMessage
		//if len(jsonMessage) < len(tempMessage) then padding of 0s is added
		//if len(jsonMessage) ? len(tempMessage) then only headerSize number of elements
		//are copied over to tempMessage
		copy(tempMessage, jsonMessage)
		//make jsonMessage equals to the new array since this will be returned
		jsonMessage = tempMessage
		//use the len(jsonMessage) and cap(jsonMessage) functions to
		//determine these values
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
