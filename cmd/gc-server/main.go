package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnmastroberti/gochat/db"
	"github.com/johnmastroberti/gochat/msg"
	homedir "github.com/mitchellh/go-homedir"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func main() {
	// Initialize database
	home, err := homedir.Dir()
	check(err)
	dbDir := filepath.Join(home, ".cache", "gochat")
	err = os.MkdirAll(dbDir, 0700)
	check(err)
	err = db.UserDBInit(filepath.Join(dbDir, "users.db"))
	check(err)

	listener, err := net.Listen("tcp", "localhost:8080")
	check(err)
	defer listener.Close()

	log.Println("Listening on port 8080")

	for {
		conn, err := listener.Accept()
		check(err)

		go handleConnection(conn)
	}

}

func handleConnection(c net.Conn) {
	log.Printf("Accepted connection from %v on %v\n", c.RemoteAddr(), c.LocalAddr())
	defer c.Close()

	client := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))

	for {
		message, err := client.ReadString('\n')
		if err != nil {
			return
		}

		log.Printf("Received message from %v: \"%s\"\n", c.RemoteAddr(), strings.Trim(message, "\n"))

		resp, dc := handleMessage(message)

		_, err = client.WriteString(resp)
		if err != nil {
			return
		}
		err = client.Flush()
		if err != nil {
			return
		}
		if dc {
			log.Println("Disconnecting from ", c.RemoteAddr())
			return
		}

	}
}

// Performs whatever actions are necessary to respond
// to the given message, and returns the response that
// should be written to the client, along with true if
// we should disconnect from the client
func handleMessage(m string) (string, bool) {
	if m == "exit\n" {
		return "See ya later\n", true
	}
	switch msg.GetMessageType([]byte(m)) {
	// New Users
	case msg.NewUserMessageT:
		newu, err := msg.NewUserFromJson([]byte(m))
		if err != nil {
			log.Println(err)
			return "disconnected\n", true
		}
		err = db.AddNewUser(newu.Username, newu.Email, newu.Password)
		if err != nil {
			log.Println(err)
			return "disconnected\n", true
		}
		return "success\n", false

		// Authenticate existing users
	case msg.LoginMessageT:
		auth, err := msg.LoginFromJson([]byte(m))
		if err != nil {
			log.Println(err)
			return "disconnected\n", true
		}
		good := db.AuthenticateUser(auth.Username, auth.Password)
		if good {
			return "success\n", false
		}
		return "disconnected\n", true

	default:
		return "Message received\n", false
	}
}
