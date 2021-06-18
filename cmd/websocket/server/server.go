package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close(websocket.StatusInternalError, "QQ")

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Get: %v \n", v)

		err = wsjson.Write(ctx, conn, "Hello")
		if err != nil {
			fmt.Println(err)
			return
		}

		conn.Close(websocket.StatusNormalClosure, "")
	})
	log.Fatal(http.ListenAndServe(":2021", nil))
}
