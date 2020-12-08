package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/johnmastroberti/gochat/msg"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// Perform the login procedure, returning true if successful
func login(stdin *bufio.Reader, server Server) bool {
	var newUser bool // true for new user, false for existing
	for {
		fmt.Println("Use '/new' to create a new user, or '/login' to login")
		line, _ := stdin.ReadString('\n')
		switch line {
		case "/new\n":
			newUser = true
			break

		case "/login\n":
			newUser = false
			break

		default:
			continue
		}
		break
	}
	var username, emailadr, password string
	// Get info
	fmt.Print("Username: ")
	username, _ = stdin.ReadString('\n')
	if newUser {
		fmt.Print("Email: ")
		emailadr, _ = stdin.ReadString('\n')
	}
	fmt.Print("Password: ")
	password, _ = stdin.ReadString('\n')
	// Send new user message to server
	if newUser {
		server.MessagesOut <- string(msg.NewUserMessage{
			Username: strings.Trim(username, "\n"),
			Email:    strings.Trim(emailadr, "\n"),
			Password: strings.Trim(password, "\n")}.ToJson())
	} else {
		server.MessagesOut <- string(msg.LoginMessage{
			Username: strings.Trim(username, "\n"),
			Password: strings.Trim(password, "\n")}.ToJson())
	}

	resp := <-server.MessagesIn
	return resp == "success\n"
}

type Server struct {
	MessagesIn  chan string
	MessagesOut chan string
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}
	defer conn.Close()

	stdin := bufio.NewReader(os.Stdin)
	serverIn := bufio.NewReader(conn)
	serverOut := bufio.NewWriter(conn)

	messagesIn := make(chan string, 100)
	messagesOut := make(chan string, 100)
	server := Server{messagesIn, messagesOut}

	go receiveMessages(serverIn, messagesIn)
	go sendMessages(serverOut, messagesOut)

	if !login(stdin, server) {
		fmt.Println("Login failed")
		return
	}

	fmt.Println("Login successful")
	go displayIncomingMessages(messagesIn)

	handleCommands(stdin, server)
}

func handleCommands(stdin *bufio.Reader, server Server) {
	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Print("Input: ", line) // TODO
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

func displayIncomingMessages(messagesIn chan string) {
	for messageString := range messagesIn {
		// Only handling TextMessages for now
		if msg.GetMessageType([]byte(messageString)) != msg.TextMessageT {
			continue
		}

		message, err := msg.TextFromJson([]byte(messageString))
		if err != nil {
			log.Println(err)
			continue
		}

		// Print the message
		fmt.Printf("New message!\nFrom: %s\nTo: %s\nContent:%s\n",
			message.From, message.To, message.Content)
	}
}
