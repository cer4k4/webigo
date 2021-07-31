package main

import(
	"log"
	"github.com/streadway/amqp"
/*	"time"
	"bytes"*/
)


func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}


func main(){
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err,"Failed to connect to RabbitMQ")
	defer conn.Close()
	ch,err := conn.Channel()
	failOnError(err,"Failed To Create Channel")
	queue,err := ch.QueueDeclare(
		"Golang-Backend",//name
		false,//durable
		false,//delete when unsend
		false,//exclusive
		false,//no-wait
		nil,//arguments
	)
	failOnError(err,"Failed To Declare queue")

	msgs,err:=ch.Consume(
		queue.Name,//queue
		"",//consumer
		true,//auto-act
		false,//clusive
		false,//no-local
		false,//no-wait
		nil,//args
	)
	failOnError(err,"Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs{
			log.Printf("Received a message: %s",d.Body)
		}
	}()
	/*go func() {
	  for d := range msgs {
	    log.Printf("Received a message: %s", d.Body)
	    dotCount := bytes.Count(d.Body, []byte("."))
	    t := time.Duration(dotCount)
	    time.Sleep(t * time.Second)
	    log.Printf("Done")
	    d.Ack(false)
	  }
	}()*/
	log.Printf("Waiting for messages.To exit preees CTRL+C")
	<-forever
}

