package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	RMQ *amqp.Connection
}

func main() {
	rmquser := os.Getenv("RABBITMQ_USER")
	rmqpass := os.Getenv("RABBITMQ_PASS")
	rmqhost := os.Getenv("RABBITMQ_HOST")
	rmqport := "5672"

	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", rmquser, rmqpass, rmqhost, rmqport)

	rmqConn, err := amqp.Dial(address)
	if err != nil {
		return
	}
	defer rmqConn.Close()

	app := &App{RMQ: rmqConn}
	defer rmqConn.Close()

	app.StartListening()

}

func (a *App) StartListening() {
	ch, err := a.RMQ.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %e", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume("messanger_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		for range msgs {
			log.Println("Message arrived, notification sending!")
			a.Notify()
		}
	}()
}

func (a *App) Notify() {

}
