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

var rabbitConnChan chan *amqp.Connection = make(chan *amqp.Connection)
var rabbitConnectionFailChan chan bool = make(chan bool)
var rabbitConsumerFailChan chan bool = make(chan bool)

func rabbitConnect() {
	addrs, err := net.LookupIP(serverConfig.RabbitAddr)
	if err != nil {
		log.Println(err)
		rabbitConnectionFailChan <- true
		return
	}

	addr := addrs[0].String() + ":" + serverConfig.RabbitPort

	conn, err := amqp.Dial("amqp://guest:guest@" + addr + "/")
	if err != nil {
		log.Println(err)
		rabbitFailChan <- true
		return
	}

	log.Println("Connected to:", addr)

	rabbitConnChan <- conn
}


func rabbitConsumer(ch *amqp.Channel, queue string, performerFunc func()) {
	q, err := ch.QueueDeclare("snowplow", true, false, false, false, nil)
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
		counter++

		log.Println(counter, "Event accepted")

		normString, nginxLog, event, err := processEventBody(d.Body)
		if err != nil {
			d.Ack(true)
			failRing.Add(normString, nginxLog, event, err.Error())
			log.Println(counter, err)
		} else {
			d.Ack(true)
			successRing.Add(normString, nginxLog, event, "success")
			log.Println(counter, "Succesfuly witten")
		}
	}
}

func rabbitChannel(conn *amqp.Connection) {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err, "Failed to open a channel")
		return
	}

	defer ch.Close()

	go rabbitConsumer(ch, "snowplow", processSnowplowEvent())
	go rabbitConsumer(ch, "adjust", processAdjustEvent())

	<- rabbitConsumerFailChan
}

func rabbitWorker() {
	for {
		go rabbitConnect()
		select {
			case conn := <- rabbitConnChan:
				rabbitChannel(conn)
				log.Println("Connection was reseted! Trying reconnect...")
			case <- rabbitConnectionFailChan:
				log.Println("Can't connect to RabbitMQ! Retry after 5s")
				time.Sleep(5*time.Second)
		}
	}
}

func processSnowplowEvent(eventBody []byte) (normString string, nginxLog *NginxLog, event *Event, err error) {
	normString, err = strconv.Unquote(`"` + string(eventBody) + `"`)
	normString = strings.Replace(normString, "\" ", "", -1)
	if err != nil { return }

	nginxLog = &NginxLog{}
	err = json.Unmarshal([]byte(normString), nginxLog)
	if err != nil { return }

	event = nginxLogToEvent(nginxLog)

	err = dbmap.Insert(event)
	if err != nil { return }

	return
}
