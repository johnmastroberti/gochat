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
func login(stdin *bufio.Reader, server *bufio.ReadWriter) bool {
	for {
		fmt.Println("Use '/new' to create a new user, or '/login' to login")
		line, _ := stdin.ReadString('\n')
		switch line {
		case "/new\n":
			// Get info
			fmt.Print("Username: ")
			username, _ := stdin.ReadString('\n')
			fmt.Print("Email: ")
			emailadr, _ := stdin.ReadString('\n')
			fmt.Print("Password: ")
			password, _ := stdin.ReadString('\n')
			// Send new user message to server
			_, err := server.Write(msg.NewUserMessage{
				Username: strings.Trim(username, "\n"),
				Email:    strings.Trim(emailadr, "\n"),
				Password: strings.Trim(password, "\n")}.ToJson())
			check(err)
			err = server.Flush()
			check(err)

			resp, err := server.ReadString('\n')
			if resp == "success\n" {
				fmt.Println("Account created successfully")
				return true
			} else {
				fmt.Println("Account creation failed")
				return false
			}

		case "/login\n":
			fmt.Print("Username: ")
			username, _ := stdin.ReadString('\n')
			fmt.Print("Password: ")
			password, _ := stdin.ReadString('\n')
			// Send new user message to server
			_, err := server.Write(msg.LoginMessage{
				Username: strings.Trim(username, "\n"),
				Password: strings.Trim(password, "\n")}.ToJson())
			check(err)
			err = server.Flush()
			check(err)

			resp, err := server.ReadString('\n')
			if resp == "success\n" {
				fmt.Println("Login successful")
				return true
			} else {
				fmt.Println("Login failed")
				return false
			}

		default:
			continue
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}
	defer conn.Close()

	stdin := bufio.NewReader(os.Stdin)
	server_in := bufio.NewReader(conn)
	server_out := bufio.NewWriter(conn)
	server := bufio.NewReadWriter(server_in, server_out)

	login(stdin, server)
	// for {
	//   fmt.Print(">")
	//   line, err := stdin.ReadString('\n')
	//   check(err)
	//
	//   server.WriteString(line)
	//   server.Flush()
	//
	//   resp, err := server.ReadString('\n')
	//   check(err)
	//   fmt.Print("Server response: ", resp)
	//
	//   if line == "exit\n" {
	//     fmt.Println("Goodbye!")
	//     return
	//   }
	//
	// }

}
