package main

import (
	"net"
	"log"
	"time"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	ConnChan     chan *amqp.Connection
	ConnFailChan chan bool
	Channel      *amqp.Channel
	ChanChan     chan *amqp.Channel
	ChanFailChan chan bool
	ConsFailChan chan bool
}

func (rabbit *Rabbit) Connect(config *Worker.Config) {
	addrs, err := net.LookupIP(config.RabbitAddr)
	if err != nil {
		log.Println(err)
		rabbit.ConnFailChan <- true
		return
	}

	addr := addrs[0].String() + ":" + config.RabbitPort

	conn, err := amqp.Dial("amqp://guest:guest@" + addr + "/")
	if err != nil {
		log.Println(err)
		rabbit.ConnFailChan <- true
		return
	}

	log.Println("Connected to:", addr)

	rabbit.ConnChan <- conn
}

func (rabbit *Rabbit) Channel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		rabbit.ChanFailChan <- true
		return
	}

	rabbit.ChanChan <- ch
}

func (rabbit *Rabbit) Consume(queueName string, bus chan []byte) {
	defer rabbit.ConsFailChan <- true

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Println(err, "Failed to declare a queue")
		return
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println(err, "Failed to consume messages")
		return
	}

	for d := range msgs {
		bus <- d.Body
		d.Ack(true)
	}
}
