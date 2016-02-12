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
	ChanChan     chan *amqp.Channel
	ChanFailChan chan bool
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
