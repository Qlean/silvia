package silvia

import (
	"net"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	ConsFailChan chan bool
}

func (rabbit *Rabbit) Connect(config *Config) error {
	addrs, err := net.LookupIP(config.RabbitAddr)
	if err != nil {
		return err
	}

	addr := addrs[0].String() + ":" + config.RabbitPort

	rabbit.Connection, err = amqp.Dial("amqp://guest:guest@" + addr + "/")
	if err != nil {
		return err
	}

	return nil
}

func (rabbit *Rabbit) Consume(queueName string, bus chan []byte) error {
	defer func() {
		rabbit.ConsFailChan <- true
	}()

	q, err := rabbit.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := rabbit.Channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for d := range msgs {
		err := d.Ack(false)
		if err != nil {
			return err
		}

		bus <- d.Body
	}

	return nil
}
