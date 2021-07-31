// websockets.go
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func Chat(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
}
func Echo(w http.ResponseWriter, r *http.Request) {
	RMQconn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
	defer RMQconn.Close()

	ch, _ := RMQconn.Channel()
	defer ch.Close()

        conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	fmt.Printf("%s %s ",conn.RemoteAddr(),"Another Client Connected\n")
	for{
            // Read message from browser
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }

            // Print the message to the console
            fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	    queue, _  := ch.QueueDeclare(
		    "Golang-Backend",//name
		    false,//durable
		    false,//delete when unsend
		    false,//exclusive
		    false,//no-wait
		    nil,//arguments
	    )

	    ch.Publish(
	    "",//exchange
	    queue.Name,//routingkey
	    false,//mandatory
	    false,//immediate
	    amqp.Publishing{
		    ContentType:"text/plain",
		    Body:  []byte(msg),
	    })

            // Write message back to browser
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
    }
}

func main() {
    http.HandleFunc("/echo",Echo)
    http.HandleFunc("/",Chat)
    http.ListenAndServe(":8080", nil)
}
