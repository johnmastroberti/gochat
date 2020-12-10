// This package deals with the messages that are
// passed between the client and server
package msg

import (
	"encoding/json"
	"strings"
)

//var ErrBadFormat = errors.New("Bad message format")

// Message is the primary data structure describing
// client/server messages
type Message struct {
	// The type of message
	Type string `json:"type,omitempty"`
	// The user to whom the message is directed
	To string `json:"to,omitempty"`
	// The user that the message is from
	From string `json:"from,omitempty"`
	// The content of the message
	Content string `json:"content,omitempty"`
	// The username to be used for login
	Username string `json:"username,omitempty"`
	// The user's email address (used to create a new account)
	Email string `json:"email,omitempty"`
	// The user's password (for login)
	Password string `json:"password,omitempty"`
}

// JSON encoding

// Encode a message as json for transfer over the network
func (m Message) ToJson() string {
	jsonText, _ := json.Marshal(m)
	return string(append(jsonText, '\n'))
}

// Convert json to a TextMessage
func FromJson(message string) (Message, error) {
	var m Message
	err := json.Unmarshal([]byte(message), &m)
	return m, err
}

// Command Parsing

// Returns true if m is a command from the user (such as "/login")
func (m Message) IsCommand() bool {
	return m.Type == "USERINPUT" && strings.IndexByte(m.Content, '/') == 0
}

// Returns the command portion of a command message (e.g. "login" for "/login")
// Returns "" if the message is not a command
func (m Message) Command() string {
	words := strings.Split(m.Content, " ")
	if len(words) > 0 {
		return strings.TrimPrefix(words[0], "/")
	}
	return ""
}
