package main

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	HostName   string
	PortNumber int
	Action     int
	Files      []File
}

type Err struct {
	HostName   string
	PortNumber int
	Code       int
}

func encodeMessage(hostName string, portNumber int, action int, files []File) []byte {
	message := Message{
		HostName:   hostName,
		PortNumber: portNumber,
		Action:     action,
		Files:      files,
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
	jsonMessage = bytes.Trim(jsonMessage, "\x00")

	err := json.Unmarshal(jsonMessage, &message)
	checkError(err)

	return message
}

func encodeError(code int) []byte {
	errObject := Err{
		HostName:   localPeer.host,
		PortNumber: localPeer.port,
		Code:       code,
	}
	errMessage, err := json.Marshal(errObject)
	checkError(err)

	return errMessage
}

const (
	Join     = iota
	Leave    = iota
	Query    = iota
	Insert   = iota
	Add      = iota
	Remove   = iota
	Files    = iota
	Upload   = iota
	Download = iota
	Have     = iota
)
