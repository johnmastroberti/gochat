// This package deals with the messages that are
// passed between the client and server
package msg

// Message is the primary data structure describing
// client/server messages
type Message struct {
	// The user to whom the message is directed
	To string `json:"to"`
	// The user that the message is from
	From string `json:"from"`
	// The content of the message
	Content string `json:"content"`
}
