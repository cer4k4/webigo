// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"fmt"
	"github.com/gorilla/websocket"
	"encoding/json"

)

type WSMessage struct {
	Sender         string `json:"sender"`
	Receiver       string `json:"receiver"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	ChatRoomID     int    `json:"chat_room_id"`
	Type           string `json:"type"`
	ChatroomCreate int    `json:"chatroom_create"`
}

var addr = flag.String("addr", "185.211.59.213:1234", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api/v1/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
	fmt.Println("Benevis")
	var txt WSMessage
	fmt.Scan(&txt.Sender)
	fmt.Scan(&txt.Receiver)
	fmt.Scan(&txt.Type)
	fmt.Scan(&txt.Message)
	fmt.Scan(&txt.ChatRoomID)
	data := map[string]WSMessage{
		"data": txt,
	}
	jtxt,_ := json.Marshal(data)
	//var parsejson WSMessage
	//_ = json.Unmarshal(jtxt,&parsejson)
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage,jtxt)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
//			case <-time.After(time.Second):
	//		case <-
			}
			return
		}
	}
}
