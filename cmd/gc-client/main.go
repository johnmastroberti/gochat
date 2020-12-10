package main

import (
	"bufio"
	"log"
	"strings"

	"github.com/johnmastroberti/gochat/msg"
	"github.com/johnmastroberti/gochat/ui"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// Perform the login procedure, returning the new user/login message to be
// sent to the server
func promptLogin(userInput chan msg.Message, uiEvents chan msg.Message, newUser bool) msg.Message {
	// Get info
	var username, emailadr, password string
	uiEvents <- msg.Message{Type: "PROMPT", Content: "Username:"}
	username = (<-userInput).Content
	if newUser {
		uiEvents <- msg.Message{Type: "PROMPT", Content: "Email:"}
		emailadr = (<-userInput).Content
	}
	uiEvents <- msg.Message{Type: "PROMPT", Content: "Password:"}
	password = (<-userInput).Content
	// Send new user message to server
	if newUser {
		return msg.Message{
			Type:     "NEWU",
			Username: strings.Trim(username, "\n"),
			Email:    strings.Trim(emailadr, "\n"),
			Password: strings.Trim(password, "\n")}
	} else {
		return msg.Message{
			Type:     "AUTH",
			Username: strings.Trim(username, "\n"),
			Password: strings.Trim(password, "\n")}
	}
}

type Server struct {
	MessagesIn  chan string
	MessagesOut chan string
}

func main() {
	//conn, err := net.Dial("tcp", "localhost:8080")
	//if err != nil {
	//	fmt.Println("Could not connect to server")
	//	return
	//}
	//defer conn.Close()

	uiEvents := make(chan msg.Message, 20)
	userInput := make(chan msg.Message, 20)
	go handleUserInput(userInput, uiEvents)
	uiEvents <- msg.Message{Type: "TEXT",
		Content: `Welcome to GoChat!
		Use /new to create a new account or /login to login`}
	ui.RunUILoop(uiEvents, userInput)

	//stdin := bufio.NewReader(os.Stdin)
	//serverIn := bufio.NewReader(conn)
	//serverOut := bufio.NewWriter(conn)

	//messagesIn := make(chan string, 100)
	//messagesOut := make(chan string, 100)
	//server := Server{messagesIn, messagesOut}

	//go receiveMessages(serverIn, messagesIn)
	//go sendMessages(serverOut, messagesOut)

	//if !login(stdin, server) {
	//	fmt.Println("Login failed")
	//	return
	//}

	//fmt.Println("Login successful")
	//go displayIncomingMessages(messagesIn)

	//handleCommands(stdin, server)
}

// Parses input from the user and sends necessary uiEvents on the uiEvents channel
func handleUserInput(userInput chan msg.Message, uiEvents chan msg.Message) {
	for m := range userInput {
		if m.IsCommand() {
			switch m.Command() {
			case "new":
				e := promptLogin(userInput, uiEvents, true)
				uiEvents <- e
			case "login":
				e := promptLogin(userInput, uiEvents, false)
				uiEvents <- e
			case "exit":
				uiEvents <- msg.Message{Content: "/exit"}
			}
		} else {
			uiEvents <- msg.Message{Type: "TEXT", Content: "Received \"" + m.Content + "\""}
		}
	}
}

// Receives messages from the server and adds them to the messagesIn channel
// Closes the channel if communication with the server fails
func receiveMessages(server *bufio.Reader, messagesIn chan string) {
	for {
		message, err := server.ReadString('\n')
		if err != nil {
			log.Println(err)
			close(messagesIn)
			return
		}
		messagesIn <- message
	}
}

// Sends messages from messagesOut to the server
// Returns when messagesOut is closed
func sendMessages(server *bufio.Writer, messagesOut chan string) {
	for message := range messagesOut {
		server.WriteString(message)
		server.Flush()
	}
}
