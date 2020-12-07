package main

import (
	"log"
	"net"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	check(err)
	defer listener.Close()

	log.Println("Listening on port 8080")

	for {
		conn, err := listener.Accept()
		check(err)
		defer conn.Close()

		go handle(conn)
	}

}

func handle(c net.Conn) {
	log.Printf("Accepted connection from %v on %v\n", c.RemoteAddr(), c.LocalAddr())

	message := make([]byte, 1024)
	for {
		message_length, err := c.Read(message)
		if err != nil {
			return
		}

		if string(message[:message_length]) == "exit\n" {
			log.Println("Disconnecting from ", c.RemoteAddr())
			c.Write([]byte("See ya later\n"))
			return
		}

		log.Printf("Received message from %v: \"%s\"\n", c.RemoteAddr(), string(message[:message_length-1]))
		c.Write([]byte("Message received\n"))
	}
}
