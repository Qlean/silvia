package main

import (
	"log"
	"strconv"
	"encoding/json"

	"github.com/streadway/amqp"
)

func rabbitConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq.service.consul/")
	if err != nil {
		log.Fatal(err, "Failed to connect to RabbitMQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "Failed to open a channel")
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"snowplow", // name
		true,   // durable
		false,    // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	for d := range msgs {
		counter++
		localCounter := counter

		var err error
		var nginxLog *NginxLog = &NginxLog{}
		var event *Event

		log.Println(counter, "Event accepted")

		// Remove Line Feed
		d.Body = d.Body[0:len(d.Body)-1]
		rawBody := string(d.Body)

		normString, err := strconv.Unquote(`"` + rawBody + `"`)
		if err != nil {
			eventFail(localCounter, "Failed to normalize nginx log string", err, rawBody, nginxLog, event)
			return
		}

		err = json.Unmarshal([]byte(normString), nginxLog)
		if err != nil {
			eventFail(localCounter, "Failed to unmarshal JSON to NginxLog structure", err, normString, nginxLog, event)
			return
		}

		d.Ack(true)

		eventsChannel <- nginxLogToEvent(nginxLog)
	}
}
