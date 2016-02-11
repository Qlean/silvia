package main

import(
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/abh/geoip"
)

type Config struct {
	Port       string
	PgConnect  string `yaml:"pg_connect"`
	RingSize   int    `yaml:"ring_size"`
	RabbitAddr string `yaml:"rabbit_addr"`
	RabbitPort string `yaml:"rabbit_port"`
}

type Worker struct {
	Config   *Config
	Stats    *Stats
	AdjustRequestBus chan *AdjustRequest
	SnowplowRequestBus chan *SnowplowRequest
	AdjustEventBus chan *AdjustEvent
	SnowplowEventBus chan *SnowplowEvent
}

func (config *Config) load() error {
	dbmap = initDb()

	defer dbmap.Db.Close()

	startTime = time.Now()

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(buf, &server_config)
	if err != nil {
		return
	}

	geo, err = geoip.Open("GeoLiteCity.dat")
	if err != nil {
		return
	}

	return
}

func (worker *Worker) generator() {
	for {
		go rabbit.connect()
		select {
			case conn := <- rabbitConnChan:
				rabbitChannel(conn)
				log.Println("Connection was reseted! Trying reconnect.")
			case <- rabbitConnectionFailChan:
				log.Println("Can't connect to RabbitMQ! Retry after 5s")
				time.Sleep(5*time.Second)
		}
	}
}

func (worker *Worker) transformer() {
	for {
		select {
			case adjustRequest := <- AdjustRequestBus:
				adjustEvent, err := adjustRequest.Transform()
				if err != nil {
					adjustEvent.Fail(err)
				}

				AdjustEventBus <- adjustEvent
			case snowplowRequest := <- SnowplowRequestBus:
				snowplowEvent, err := snowplowRequest.Transform()
				if err != nil {
					snowplowEvent.Fail(err)
				}

				AdjustEventBus <- adjustEvent
		}
	}
}

func (worker *Worker) capacitor() {
	for {
		go postgresConnect()
		select {
			case conn := <- postgresConnChan:
				postgresWriter(conn)
				log.Println("PostgreSQL connection was reseted! Trying reconnect.")
			case <- postgresConnectionFailChan:
				log.Println("Can't connect to PostgreSQL! Retry after 5s")
				time.Sleep(5*time.Second)
		}
	}
}

