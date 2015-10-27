package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	var c chan string = make(chan string)
	go checker(c)
	go publish(ch, q, c)
	var input string
	fmt.Scanln(&input)
}

func publish(ch *amqp.Channel, q amqp.Queue, c <- chan string) {
	for {
		connectionStatus := <- c
		if connectionStatus == "lost" {
			for i := 0; i < 40 ; i++ {
				conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
				if err != nil {
					fmt.Println("Failed to connect to RabbitMQ")
					time.Sleep(1 * time.Second)
					continue
				} else {
					defer conn.Close()
				}
			}
		}
		if connectionStatus == "ok" {
			conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
			if err != nil {
				fmt.Println("Failed to connect to RabbitMQ")
				time.Sleep(1 * time.Second)
				continue
			} else {
				defer conn.Close()
			}
			ch, _ := conn.Channel()
			defer ch.Close()
			body := bodyFrom(os.Args)
			q, err := ch.QueueDeclare(
				"task_queue", // name
				true,         // durable
				false,        // delete when unused
				false,        // exclusive
				false,        // no-wait
				nil,          // arguments
			)
			err = ch.Publish(
				"",           // exchange
				q.Name,       // routing key
				false,        // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
				})
			if err == nil {
				log.Printf(" [x] Sent %s", body)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func checker(c chan <- string) {
	for {
		_, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			c <- "lost"
		} else {
			c <- "ok"
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}