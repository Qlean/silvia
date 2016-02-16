package main

import(
	"time"
	"reflect"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/abh/geoip"

	consul "github.com/hashicorp/consul/api"
)

type(
	Stats struct {
		AdjustSuccessRing   *AdjustRing
		AdjustFailRing      *AdjustRing
		SnowplowSuccessRing *SnowplowRing
		SnowplowFailRing    *SnowplowRing
		StartTime           time.Time
		RabbitHealth        bool
		PostgresHealth      bool
	}

	Config struct {
		Port       string `consul:"port"`
		PgConnect  string `consul:"pg_connect"`
		RingSize   string `consul:"ring_size"`
		RabbitAddr string `consul:"rabbit_addr"`
		RabbitPort string `consul:"rabbit_port"`
	}

	Worker struct {
		Config             *Config
		Stats              *Stats
		AdjustRequestBus   chan []byte
		SnowplowRequestBus chan []byte
		AdjustEventBus     chan *AdjustEvent
		SnowplowEventBus   chan *SnowplowEvent
		GeoDB              *geoip.GeoIP
	}
)

func (config *Config) FillFromConsul(client *consul.Client, appName string) error {
	structType := reflect.TypeOf(*config)
	structValue := reflect.ValueOf(config).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		consulKey := field.Tag.Get("consul")
		pair, _, _ := kv.Get(appName + "/" + consulKey, nil)
		if pair == nil {
			return errors.New("Failed to fetch - " + consulKey)
		}
		consulValue = string(pair.Value)
		structValue.FieldByIndex(i).SetString(consulValue)
	}

	return nil
}


func (worker *Worker) load() error {
	stats := worker.Stats

	stats.StartTime = time.Now()
	stats.SuccessRing = &Ring{}
	stats.FailRing = &Ring{}

	geo, err = geoip.Open("GeoLiteCity.dat")
	if err != nil {
		return err
	}

	consulConfig := &consul.Config{
		Address:    "127.0.0.1:8500",
		Scheme:     "http",
		HttpClient: http.DefaultClient,
	}

	client, err := consul.NewClient(consulConfig)
	if err != nil {
		return err
	}

	err := worker.Config.FillFromConsul(client, "silvia")
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) generator() {
	for {
		rabbit := &Rabbit{}
		go rabbit.Connect(worker.Config)
		select {
			case conn := <- rabbit.ConnChan:
				go rabbit.NewChannel(conn)
				select {
					case rabbit.Channel := <- rabbit.ChanChan:
						worker.Stats.RabbitHealth = true
						go rabbit.Consume("adjust", worker.AdjustRequestBus)
						go rabbit.Consume("snowplow", worker.SnowplowRequestBus)
						<- rabbit.ConsFailChan
						rabbit.Channel.Close()
					case <- rabbit.ChannelFailChan:
						log.Println("Can't create RabbitMQ channel")
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
	go func() {
		for {
			adjustRequest := <- worker.AdjustRequestBus:
			adjustEvent := &AdjustEvent{}
			err := adjustEvent.Transform(adjustRequest)
			if err != nil {
				worker.Stats.AdjustFailRing.Add(adjustEvent, err)
			} else {
				worker.AdjustEventBus <- adjustEvent
			}
		}
	}()

	go func() {
		for {
			snowplowRequest := <- worker.SnowplowRequestBus:
			snowplowEvent := &SnowplowEvent{}
			err := snowplowEvent.Transform(snowplowRequest)
			if err != nil {
				worker.Stats.SnowplowFailRing.Add(snowplowEvent, err)
			} else {
				worker.SnowplowEventBus <- snowplowEvent
			}
		}
	}()
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
								worker.Stats.AdjustFailRing.Add(adjustEvent, err)
							} else {
								worker.Stats.AdjustSuccessRing.Add(adjustEvent, err)
							}
						case snowplowEvent := <- worker.SnowplowEventBus:
							err := snowplowEvent.Write()
							if err != nil {
								worker.Stats.SnowplowFailRing.Add(snowplowEvent, err)
							} else {
								worker.Stats.SnowplowSuccessRing.Add(snowplowEvent, err)
							}
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
