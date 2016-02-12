package main

import(
	"sync"
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/abh/geoip"
)

type(
	Stats struct {
		SuccessRing    *Ring
		FailRing       *Ring
		StartTime      time.Time
		RabbitHealth   bool
		PostgresHealth bool
	}

	Config struct {
		Port       string `consul:"port"`
		PgConnect  string `consul:"pg_connect"`
		RingSize   int    `consul:"ring_size"`
		RabbitAddr string `consul:"rabbit_addr"`
		RabbitPort string `consul:"rabbit_port"`
	}

	Worker struct {
		Config   *Config
		Stats    *Stats
		AdjustRequestBus chan *AdjustRequest
		SnowplowRequestBus chan *SnowplowRequest
		AdjustEventBus chan *AdjustEvent
		SnowplowEventBus chan *SnowplowEvent
	}

	Consumer interface {
		Consume()
	}

	Transformer interface {
		Transform()
	}

	Writer interface {
		Write()
	}
)


func (config *Config) load() error {
	startTime = time.Now()


	geo, err = geoip.Open("GeoLiteCity.dat")
	if err != nil {
		return
	}

	return
}

func (worker *Worker) generator() {
	for {
		rabbit := &Rabbit{}
		go rabbit.Connect(worker.Config)
		select {
			case conn := <- rabbit.ConnChan:
				go rabbit.NewChannel(conn)
				select {
					case ch := <- rabbit.ChannelChan:
						worker.Stats.RabbitHealth = true
						wg := sync.WaitGroup()
							go worker.AdjustRequestBus.Consume(ch, wg)
							go worker.SnowplowRequestBus.Consume(ch, wg)
						wg.Wait()
						ch.Close()
					case <- rabbit.ChannelFailChan:
						log.Println("Can't create RabbitMQ channel.")
				}
				conn.Close()
				log.Println("Connection was reseted! Retry after 5s")
			case <- rabbit.ConnectionFailChan:
				log.Println("Can't connect to RabbitMQ! Retry after 5s")
		}
		worker.Stats.RabbitHealth = false
		time.Sleep(5*time.Second)
	}
}

func (worker *Worker) transformer() {
	for {
		select {
			case adjustRequest := <- worker.AdjustRequestBus:
				adjustEvent, err := adjustRequest.Transform()
				if err != nil {
					adjustEvent.Fail(err)
				}

				worker.AdjustEventBus <- adjustEvent
			case snowplowRequest := <- worker.SnowplowRequestBus:
				snowplowEvent, err := snowplowRequest.Transform()
				if err != nil {
					snowplowEvent.Fail(err)
				}

				worker.AdjustEventBus <- adjustEvent
		}
	}
}

func (worker *Worker) capacitor() {
	for {
		postgres := &Postgres{}
		go postgres.Connect(worker.Config)
		select {
			case conn := <- postgres.ConnChan:
				worker.Stats.PostgresHealth = true
				for {
					select {
						case adjustEvent := <- worker.AdjustEventBus:
							err := adjustEvent.Write()
							if err != nil {
								adjustEvent.Fail(err)
							}
							adjustEvent.Success()
						case snowplowEvent := <- worker.SnowplowEventBus:
							err := snowplowEvent.Write()
							if err != nil {
								snowplowEvent.Fail(err)
							}
							snowplowEvent.Success()
					}
				}
				conn.Db.Close()
				log.Println("PostgreSQL connection was reseted! Retry after 5s")
			case <- postgres.ConnFailChan:
				log.Println("Can't connect to PostgreSQL! Retry after 5s")
		}
		worker.Stats.PostgresHealth = false
		time.Sleep(5*time.Second)
	}
}

