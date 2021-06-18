package main

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil) // ws -> http automatically in Dial
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "QQ")

	err = wsjson.Write(ctx, c, "YO1")
	if err != nil {
		panic(err)
	}
	err = wsjson.Write(ctx, c, "YO2")
	if err != nil {
		panic(err)
	}
	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client get: %v \n", v)

	c.Close(websocket.StatusNormalClosure, "")
}
