package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan string, 8)
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	// go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	user := &User{
		ID:             getUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}
	defer close(user.MessageChannel)
	go sendMessage(conn, user.MessageChannel)

	user.MessageChannel <- "Welcome, " + strconv.Itoa(user.ID)
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has entered"

	enteringChannel <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
	}

	if err := input.Err(); err != nil {
		fmt.Println("Read error", err)
	}

	leavingChannel <- user
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has left"
}

func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)

		case msg := <-messageChannel:
			for user := range users {
				user.MessageChannel <- msg
			}

		}
	}
}
