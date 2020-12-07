package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
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
	for {
		fmt.Print(">")
		line, err := stdin.ReadString('\n')
		check(err)

		server.WriteString(line)
		server.Flush()

		resp, err := server.ReadString('\n')
		check(err)
		fmt.Print("Server response: ", resp)

		if line == "exit\n" {
			fmt.Println("Goodbye!")
			return
		}

	}

}
