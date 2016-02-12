package main

import(
	"sync"
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/abh/geoip"
)

type(
	Config struct {
		Port       string
		PgConnect  string `yaml:"pg_connect"`
		RingSize   int    `yaml:"ring_size"`
		RabbitAddr string `yaml:"rabbit_addr"`
		RabbitPort string `yaml:"rabbit_port"`
	}

	Worker struct {
		Config   *Config
		Stats    *Stats
		AdjustRequestBus chan *AdjustRequest
		SnowplowRequestBus chan *SnowplowRequest
		AdjustEventBus chan *AdjustEvent
		SnowplowEventBus chan *SnowplowEvent
	}
)


func (config *Config) load() error {
	dbmap = initDb()

	defer dbmap.Db.Close()

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
		time.Sleep(5*time.Second)
	}
}

