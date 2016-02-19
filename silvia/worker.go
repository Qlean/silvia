package silvia

import(
	"log"
	"time"
	"strconv"
	"reflect"
	"net/http"
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

func (config *Config) fillFromConsul(client *consul.Client, appName string) error {
	kv := client.KV()

	structType := reflect.TypeOf(*config)
	structValue := reflect.ValueOf(config).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		consulKey := field.Tag.Get("consul")
		pair, _, err := kv.Get(appName + "/" + consulKey, nil)
		if pair == nil {
			return err
		}
		consulValue := string(pair.Value)
		structValue.FieldByName(field.Name).SetString(consulValue)
	}

	return nil
}


func (worker *Worker) Load() error {
	worker.Config = &Config{}

	var err error
	worker.GeoDB, err = geoip.Open("GeoLiteCity.dat")
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

	err = worker.Config.fillFromConsul(client, "silvia")
	if err != nil {
		return err
	}

	ringSize, err := strconv.Atoi(worker.Config.RingSize)
	if err != nil {
		return err
	}

	worker.Stats = &Stats{
		AdjustSuccessRing: &AdjustRing{Size: ringSize},
		AdjustFailRing: &AdjustRing{Size: ringSize},
		SnowplowSuccessRing: &SnowplowRing{Size: ringSize},
		SnowplowFailRing: &SnowplowRing{Size: ringSize},
	}

	worker.Stats.StartTime = time.Now()

	worker.AdjustRequestBus = make(chan []byte)
	worker.SnowplowRequestBus = make(chan []byte)
	worker.AdjustEventBus = make(chan *AdjustEvent)
	worker.SnowplowEventBus = make(chan *SnowplowEvent)

	check := &consul.AgentServiceCheck{
		HTTP: "http://localhost:" + worker.Config.Port + "/v1/status",
		Interval: "10s",
		Timeout: "1s",
	}

	port, err := strconv.Atoi(worker.Config.Port)
	if err != nil {
		return nil
	}

	service := &consul.AgentServiceRegistration{
		Name: "silvia",
		Port: port,
		Check: check,
	}

	agent := client.Agent()
	err = agent.ServiceRegister(service)
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) Generator() {
	for {
		rabbit := &Rabbit{}
		err := rabbit.Connect(worker.Config)
		if err != nil {
			log.Println("Can't connect to RabbitMQ! Retry after 5s")
		} else {
			rabbit.Channel, err = rabbit.Connection.Channel()
			if err != nil {
				log.Println("Can't create RabbitMQ channel! Retry after 5s")
			} else {
				worker.Stats.RabbitHealth = true
				rabbit.ConsFailChan = make(chan bool)
				go rabbit.Consume("adjust", worker.AdjustRequestBus)
				go rabbit.Consume("snowplow", worker.SnowplowRequestBus)
				<- rabbit.ConsFailChan
				rabbit.Channel.Close()
			}
			rabbit.Connection.Close()
		}

		worker.Stats.RabbitHealth = false
		time.Sleep(5*time.Second)
	}
}

func (worker *Worker) Transformer() {
	go func() {
		for {
			adjustEvent := &AdjustEvent{}
			err := adjustEvent.Transform(<- worker.AdjustRequestBus)
			if err != nil {
				worker.Stats.AdjustFailRing.Add(adjustEvent, err)
			} else {
				worker.AdjustEventBus <- adjustEvent
			}
		}
	}()

	go func() {
		for {
			snowplowEvent := &SnowplowEvent{}
			err := snowplowEvent.Transform(<- worker.SnowplowRequestBus, worker.GeoDB)
			if err != nil {
				worker.Stats.SnowplowFailRing.Add(snowplowEvent, err)
			} else {
				worker.SnowplowEventBus <- snowplowEvent
			}
		}
	}()
}

func (worker *Worker) Writer() {
	for {
		postgres := &Postgres{}
		err := postgres.Connect(worker.Config)
		if err != nil {
			log.Println("Can't connect to PostgreSQL! Retry after 5s")
		} else {
			worker.Stats.PostgresHealth = true
			for {
				select {
				case adjustEvent := <- worker.AdjustEventBus:
					err := postgres.Connection.Insert(adjustEvent)
					if err != nil {
						worker.Stats.AdjustFailRing.Add(adjustEvent, err)
					} else {
						worker.Stats.AdjustSuccessRing.Add(adjustEvent, err)
					}
				case snowplowEvent := <- worker.SnowplowEventBus:
					err := postgres.Connection.Insert(snowplowEvent)
					if err != nil {
						worker.Stats.SnowplowFailRing.Add(snowplowEvent, err)
					} else {
						worker.Stats.SnowplowSuccessRing.Add(snowplowEvent, err)
					}
				}
			}
			postgres.Connection.Db.Close()
		}

		worker.Stats.PostgresHealth = false
		time.Sleep(5*time.Second)
	}
}
