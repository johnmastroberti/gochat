// This package deals with the messages that are
// passed between the client and server
package msg

import (
	"encoding/json"
)

//var ErrBadFormat = errors.New("Bad message format")

type MessageType int8

const (
	TextMessageT MessageType = iota
	NewUserMessageT
	LoginMessageT
	BadFormatT
)

// Message is the primary data structure describing
// client/server messages
type TextMessage struct {
	// The user to whom the message is directed
	To string `json:"to"`
	// The user that the message is from
	From string `json:"from"`
	// The content of the message
	Content string `json:"content"`
}

// Encode a standard message as json
func (m TextMessage) ToJson() []byte {
	json, _ := json.Marshal(m)
	json = append(json, '\n')
	return append([]byte("TEXT:"), json...)
}

// The message that a client sends the server to create a new user
type NewUserMessage struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Encode a new user message as json
func (m NewUserMessage) ToJson() []byte {
	json, _ := json.Marshal(m)
	json = append(json, '\n')
	return append([]byte("NEWU:"), json...)
}

// The message that a client sends the server to create login
type LoginMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Encode a new user message as json
func (m LoginMessage) ToJson() []byte {
	json, _ := json.Marshal(m)
	json = append(json, '\n')
	return append([]byte("AUTH:"), json...)
}

// Determine the type of message received
func GetMessageType(message []byte) MessageType {
	if len(message) < 5 {
		return BadFormatT
	}
	key := string(message[:4])
	switch key {
	case "TEXT":
		return TextMessageT
	case "NEWU":
		return NewUserMessageT
	case "AUTH":
		return LoginMessageT
	default:
		return BadFormatT
	}
}

// Convert json to a TextMessage
func TextFromJson(message []byte) (TextMessage, error) {
	var m TextMessage
	err := json.Unmarshal(message[5:], &m)
	return m, err
}

// Convert json to a NewUserMessage
func NewUserFromJson(message []byte) (NewUserMessage, error) {
	var m NewUserMessage
	err := json.Unmarshal(message[5:], &m)
	return m, err
}

// Convert json to a LoginMesssage
func LoginFromJson(message []byte) (LoginMessage, error) {
	var m LoginMessage
	err := json.Unmarshal(message[5:], &m)
	return m, err
}
