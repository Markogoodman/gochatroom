package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

var (
	enteringChannel            = make(chan *User)
	leavingChannel             = make(chan *User)
	messageChannel             = make(chan Message, 8)
	getUserID       func() int = func() func() int {
		id := -1
		return func() int {
			id++
			return id
		}
	}()
)

type Message struct {
	OwnerID int
	Content string
}

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn)
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

	user.MessageChannel <- "Welcome, " + strconv.Itoa(user.ID) + "\n"
	messageChannel <- Message{
		OwnerID: user.ID,
		Content: "user: `" + strconv.Itoa(user.ID) + "` has entered",
	}

	enteringChannel <- user

	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Second
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- Message{
			OwnerID: user.ID,
			Content: strconv.Itoa(user.ID) + ":" + input.Text(),
		}
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		fmt.Println("Read error", err)
	}

	leavingChannel <- user
	messageChannel <- Message{
		OwnerID: user.ID,
		Content: "user: `" + strconv.Itoa(user.ID) + "` has left",
	}
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
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content + "\n"
			}

		}
	}
}
