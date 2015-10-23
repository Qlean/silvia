package main

import (
	"log"
	"strconv"
	"strings"
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

	q, err := ch.QueueDeclare("snowplow", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

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

func processEventBody(eventBody []byte) (normString string, nginxLog *NginxLog, event *Event, err error) {
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
