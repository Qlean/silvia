package silvia

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/abh/geoip"
	consul "github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

type (
	Health struct {
		sync.RWMutex
		Health bool
	}

	Stats struct {
		AdjustSuccessRing   *AdjustRing
		AdjustFailRing      *AdjustRing
		SnowplowSuccessRing *SnowplowRing
		SnowplowFailRing    *SnowplowRing

		RedshiftAdjustSuccessRing   *AdjustRing
		RedshiftAdjustFailRing      *AdjustRing
		RedshiftSnowplowSuccessRing *SnowplowRing
		RedshiftSnowplowFailRing    *SnowplowRing

		PostgresAdjustSuccessRing   *AdjustRing
		PostgresAdjustFailRing      *AdjustRing
		PostgresSnowplowSuccessRing *SnowplowRing
		PostgresSnowplowFailRing    *SnowplowRing

		StartTime      time.Time
		RabbitHealth   Health
		PostgresHealth Health
		RedshiftHealth Health
	}

	Config struct {
		PostgresEnabled string `consul:"postgres_enabled"`
		RedshiftEnabled string `consul:"redshift_enabled"`
		Port            string `consul:"port"`
		PostgresConnect string `consul:"postgres_connect"`
		RedshiftConnect string `consul:"redshift_connect"`
		RingSize        string `consul:"ring_size"`
		RabbitAddr      string `consul:"rabbit_addr"`
		RabbitPort      string `consul:"rabbit_port"`
	}

	Worker struct {
		Config                   *Config
		Stats                    *Stats
		AdjustRequestBus         chan []byte
		SnowplowRequestBus       chan []byte
		AdjustErrorBus           chan *AdjustEvent
		SnowplowErrorBus         chan *SnowplowEvent
		PostgresAdjustEventBus   chan *AdjustEvent
		PostgresSnowplowEventBus chan *SnowplowEvent
		RedshiftAdjustEventBus   chan *AdjustEvent
		RedshiftSnowplowEventBus chan *SnowplowEvent
		GeoDB                    *geoip.GeoIP
		ConsulAgent              *consul.Agent
		ConsulServiceID          string
	}
)

func (health *Health) Set(value bool) {
	health.Lock()
	health.Health = value
	health.Unlock()
}

func (health *Health) Get() bool {
	health.RLock()
	defer health.RUnlock()
	return health.Health
}

func (config *Config) fillFromConsul(client *consul.Client, appName string) error {
	kv := client.KV()

	structType := reflect.TypeOf(*config)
	structValue := reflect.ValueOf(config).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		consulKey := field.Tag.Get("consul")
		pair, _, err := kv.Get(appName+"/"+consulKey, nil)
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
		AdjustSuccessRing:   &AdjustRing{Size: ringSize},
		AdjustFailRing:      &AdjustRing{Size: ringSize},
		SnowplowSuccessRing: &SnowplowRing{Size: ringSize},
		SnowplowFailRing:    &SnowplowRing{Size: ringSize},

		RedshiftAdjustSuccessRing:   &AdjustRing{Size: ringSize},
		RedshiftAdjustFailRing:      &AdjustRing{Size: ringSize},
		RedshiftSnowplowSuccessRing: &SnowplowRing{Size: ringSize},
		RedshiftSnowplowFailRing:    &SnowplowRing{Size: ringSize},

		PostgresAdjustSuccessRing:   &AdjustRing{Size: ringSize},
		PostgresAdjustFailRing:      &AdjustRing{Size: ringSize},
		PostgresSnowplowSuccessRing: &SnowplowRing{Size: ringSize},
		PostgresSnowplowFailRing:    &SnowplowRing{Size: ringSize},
	}

	worker.Stats.StartTime = time.Now()

	worker.AdjustRequestBus = make(chan []byte)
	worker.SnowplowRequestBus = make(chan []byte)
	worker.PostgresAdjustEventBus = make(chan *AdjustEvent)
	worker.PostgresSnowplowEventBus = make(chan *SnowplowEvent)
	worker.RedshiftAdjustEventBus = make(chan *AdjustEvent, 50)
	worker.RedshiftSnowplowEventBus = make(chan *SnowplowEvent, 300)

	worker.AdjustErrorBus = make(chan *AdjustEvent, 15)
	worker.SnowplowErrorBus = make(chan *SnowplowEvent, 15)

	port, err := strconv.Atoi(worker.Config.Port)
	if err != nil {
		return nil
	}

	worker.ConsulServiceID = uuid.NewV4().String()

	checks := consul.AgentServiceChecks{
		&consul.AgentServiceCheck{
			TTL: "10s",
		},
		&consul.AgentServiceCheck{
			HTTP:     "http://localhost:" + worker.Config.Port + "/v1/status",
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	service := &consul.AgentServiceRegistration{
		ID:     worker.ConsulServiceID,
		Name:   "silvia",
		Port:   port,
		Checks: checks,
	}

	worker.ConsulAgent = client.Agent()
	err = worker.ConsulAgent.ServiceRegister(service)
	if err != nil {
		return err
	}

	go func() {
		for {
			<-time.After(8 * time.Second)
			err = worker.ConsulAgent.PassTTL("service:"+worker.ConsulServiceID+":1", "Internal TTL ping")
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("Service registred with ID:", worker.ConsulServiceID)

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
				worker.Stats.RabbitHealth.Set(true)
				rabbit.ConsFailChan = make(chan bool)
				go rabbit.Consume("adjust", worker.AdjustRequestBus)
				go rabbit.Consume("snowplow", worker.SnowplowRequestBus)
				<-rabbit.ConsFailChan
				rabbit.Channel.Close()
			}
			rabbit.Connection.Close()
		}

		worker.Stats.RabbitHealth.Set(false)
		time.Sleep(5 * time.Second)
	}
}

func (worker *Worker) Transformer() {
	go func() {
		for {
			rawEvent := <-worker.AdjustRequestBus
			adjustEvent := &AdjustEvent{}
			err := adjustEvent.Transform(rawEvent)
			if err != nil {
				worker.Stats.AdjustFailRing.Add(adjustEvent, err)
				checkStringForNull("Transform", &adjustEvent.ErrType)
				checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &adjustEvent.Error)
				checkStringForNull(strings.Replace(fmt.Sprintf("%#v", string(rawEvent)), "'", "''", -1), &adjustEvent.ErrorEvent)
				worker.AdjustErrorBus <- adjustEvent
			} else {
				if worker.Stats.PostgresHealth.Get() {
					worker.PostgresAdjustEventBus <- adjustEvent
				}
				if worker.Stats.RedshiftHealth.Get() {
					worker.RedshiftAdjustEventBus <- adjustEvent
				}
			}
		}
	}()

	go func() {
		for {
			rawEvent := <-worker.SnowplowRequestBus
			snowplowEvent := &SnowplowEvent{}
			err := snowplowEvent.Transform(rawEvent, worker.GeoDB)
			if err != nil {
				worker.Stats.SnowplowFailRing.Add(snowplowEvent, err)

				checkStringForNull("error", &snowplowEvent.EventID)
				checkStringForNull("Transform", &snowplowEvent.ErrType)
				checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &snowplowEvent.Error)
				checkStringForNull(strings.Replace(fmt.Sprintf("%#v", string(rawEvent)), "'", "''", -1), &snowplowEvent.ErrorEvent)
				worker.SnowplowErrorBus <- snowplowEvent
			} else {
				if worker.Stats.PostgresHealth.Get() {
					worker.PostgresSnowplowEventBus <- snowplowEvent
				}
				if worker.Stats.RedshiftHealth.Get() {
					worker.RedshiftSnowplowEventBus <- snowplowEvent
				}
			}
		}
	}()
}

func (worker *Worker) Writer(driver string) {

	switch driver {

	case "postgres":

		if worker.Config.PostgresEnabled == "true" {

			postgres := &Postgres{}
			err := postgres.Connect(worker.Config)

			if err != nil {
				log.Println("Can't connect to PostgreSQL! Retry after 5s")
			} else {

				worker.Stats.PostgresHealth.Set(true)
				go func() {
					defer postgres.Connection.Db.Close()
					for {
						adjustEvent := <-worker.PostgresAdjustEventBus
						err := postgres.Connection.Insert(adjustEvent)
						if err != nil {
							worker.Stats.PostgresAdjustFailRing.Add(adjustEvent, err)
						} else {
							worker.Stats.PostgresAdjustSuccessRing.Add(adjustEvent, err)
						}
					}
				}()

				go func() {
					defer postgres.Connection.Db.Close()
					for {
						snowplowEvent := <-worker.PostgresSnowplowEventBus
						err := postgres.Connection.Insert(snowplowEvent)
						if err != nil {
							worker.Stats.PostgresSnowplowFailRing.Add(snowplowEvent, err)
						} else {
							worker.Stats.PostgresSnowplowSuccessRing.Add(snowplowEvent, err)
						}
					}
				}()
			}
		}
	case "redshift":

		redshift := &Redshift{}

		if worker.Config.RedshiftEnabled == "true" {
			err := redshift.Connect(worker.Config)
			if err != nil {
				log.Println("Can't connect to Redshift! Retry after 5s")
			} else {

				// Create tablemaps for adjust and atomic
				snowplowTmap := redshift.Connection.AddTableWithNameAndSchema(SnowplowEvent{}, "atomic", "events")
				adjustTmap := redshift.Connection.AddTableWithNameAndSchema(AdjustEvent{}, "adjust", "events")

				snowplowInsert := fmt.Sprintf("INSERT INTO \"%s\".\"%s\" ( \"%s\" ) values ", "atomic", "events", strings.Join(GetColumns(snowplowTmap)[2:], "\", \""))
				adjustInsert := fmt.Sprintf("INSERT INTO \"%s\".\"%s\" (\"%s\") values ", "adjust", "events", strings.Join(GetColumns(adjustTmap)[2:], "\", \""))

				worker.Stats.RedshiftHealth.Set(true)
				go func() {
					defer redshift.Connection.Db.Close()
					for {
						var query bytes.Buffer
						query.WriteString(adjustInsert)

						remains := 5
						i := 0
						for event := range worker.RedshiftAdjustEventBus {
							i++
							// event := <-worker.RedshiftAdjustEventBus
							stringEvent, err := getStringEventValues(event)

							if err != nil {
								worker.Stats.RedshiftAdjustFailRing.Add(event, err)
								checkStringForNull("GetStringEventValues", &event.ErrType)
								checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
								checkStringForNull(strings.Replace(fmt.Sprintf("%#v", event), "'", "''", -1), &event.ErrorEvent)
								worker.AdjustErrorBus <- event
								continue
							}

							_, err = query.WriteString(stringEvent)

							if err != nil {
								worker.Stats.RedshiftAdjustFailRing.Add(event, err)
								checkStringForNull("WriteString", &event.ErrType)
								checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
								checkStringForNull(strings.Replace(fmt.Sprintf("%#v", event), "'", "''", -1), &event.ErrorEvent)
								worker.AdjustErrorBus <- event
								break
							}
							worker.Stats.RedshiftAdjustSuccessRing.Add(event, err)

							if i == remains {
								break
							}
							query.WriteString(", ")

						}

						query.WriteString(";")
						_, err = redshift.Connection.Exec(query.String())

						if err != nil {
							worker.Stats.RedshiftAdjustFailRing.Add(&AdjustEvent{}, err)
							event := &AdjustEvent{}

							checkStringForNull("ExecQuery", &event.ErrType)
							checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
							checkStringForNull(strings.Replace(fmt.Sprintf("%#v", query.String()), "'", "''", -1), &event.ErrorEvent)
							worker.AdjustErrorBus <- event
							fmt.Println("ERROR ", err, " ", query.String())
						} else {
							worker.Stats.RedshiftAdjustSuccessRing.Add(&AdjustEvent{}, err)
						}
					}
				}()

				go func() {
					defer redshift.Connection.Db.Close()
					for {
						var query bytes.Buffer
						query.WriteString(adjustInsert)

						remains := 1
						i := 0
						for event := range worker.AdjustErrorBus {
							i++
							// event := <-worker.RedshiftAdjustEventBus
							stringEvent, err := getStringEventValues(event)

							if err != nil {
								log.Println(event)
								log.Println(err)
								continue
							}

							_, err = query.WriteString(stringEvent)

							if err != nil {
								log.Println(stringEvent)
								log.Println(err)
								break
							}

							if i == remains {
								break
							}
							query.WriteString(", ")

						}

						query.WriteString(";")
						_, err = redshift.Connection.Exec(query.String())
						if err != nil {
							log.Println(query.String())
							log.Println(err)
						}
					}
				}()

				go func() {
					defer redshift.Connection.Db.Close()
					for {
						var query bytes.Buffer
						query.WriteString(snowplowInsert)
						i := 0
						remains := 50
						for event := range worker.RedshiftSnowplowEventBus {
							i++
							// for i := 0; i < remains; i++ {
							// event := <-worker.RedshiftSnowplowEventBus
							stringEvent, err := getStringEventValues(event)

							if err != nil {
								worker.Stats.RedshiftSnowplowFailRing.Add(event, err)
								checkStringForNull("error", &event.EventID)
								checkStringForNull("GetStringEventValues", &event.ErrType)
								checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
								checkStringForNull(strings.Replace(fmt.Sprintf("%#v", event), "'", "''", -1), &event.ErrorEvent)
								worker.SnowplowErrorBus <- event
								continue
							}

							_, err = query.WriteString(stringEvent)

							if err != nil {
								worker.Stats.RedshiftSnowplowFailRing.Add(event, err)
								checkStringForNull("error", &event.EventID)
								checkStringForNull("WriteString", &event.ErrType)
								checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
								checkStringForNull(strings.Replace(fmt.Sprintf("%#v", event), "'", "''", -1), &event.ErrorEvent)
								worker.SnowplowErrorBus <- event
								break
							}
							worker.Stats.RedshiftSnowplowSuccessRing.Add(event, err)

							if i == remains {
								break
							}
							query.WriteString(", ")
						}

						query.WriteString(";")

						_, err = redshift.Connection.Exec(query.String())

						if err != nil {
							event := &SnowplowEvent{}
							worker.Stats.RedshiftSnowplowFailRing.Add(event, err)

							checkStringForNull("error", &event.EventID)
							checkStringForNull("ExecQuery", &event.ErrType)
							checkStringForNull(strings.Replace(err.Error(), "'", "''", -1), &event.Error)
							checkStringForNull(strings.Replace(fmt.Sprintf("%#v", query.String()), "'", "''", -1), &event.ErrorEvent)
							worker.SnowplowErrorBus <- event
						} else {
							worker.Stats.RedshiftSnowplowSuccessRing.Add(&SnowplowEvent{}, err)
						}
					}
				}()

				go func() {
					defer redshift.Connection.Db.Close()
					for {
						var query bytes.Buffer
						query.WriteString(snowplowInsert)

						remains := 1
						i := 0
						for event := range worker.SnowplowErrorBus {
							i++
							// event := <-worker.RedshiftAdjustEventBus
							stringEvent, err := getStringEventValues(event)

							if err != nil {
								log.Println(event)
								log.Println(err)
								continue
							}

							_, err = query.WriteString(stringEvent)

							if err != nil {
								log.Println(stringEvent)
								log.Println(err)
								break
							}

							if i == remains {
								break
							}
							query.WriteString(", ")

						}

						query.WriteString(";")
						_, err = redshift.Connection.Exec(query.String())
						if err != nil {
							log.Println(query.String())
							log.Println(err)
						}
					}
				}()
			}

		}

	}

	// worker.Stats.PostgresHealth.Set(false)
	time.Sleep(3 * time.Second)
}

func (worker *Worker) Killer() {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	worker.ConsulAgent.ServiceDeregister(worker.ConsulServiceID)
	os.Exit(0)
}

func getEventValues(event interface{}) []interface{} {
	var values []interface{}
	e := reflect.ValueOf(event).Elem()
	for i := 2; i < e.NumField(); i++ {
		// if strings.Contains(e.Type().Field(i).Type.String(), "sql.") {
		// 	values = append(values, e.Field(i).Field(0).Interface())
		// 	continue
		// }
		values = append(values, e.Field(i).Interface())
	}
	return values
}

func getStringEventValues(event interface{}) (string, error) {

	var values bytes.Buffer

	e := reflect.ValueOf(event).Elem()
	values.WriteString("( ")

	formatString := "'%v', "

	for i := 2; i < e.NumField(); i++ {

		if i == e.NumField()-1 {
			formatString = "'%v' )"
		}
		// fmt.Println(e.Field(i).Type().String())
		if e.Field(i).Type().String() == "silvia.NullTime" {
			time := e.Field(i).Field(0).MethodByName("Format").Call([]reflect.Value{reflect.ValueOf(time.RFC3339Nano)})[0]
			values.WriteString(fmt.Sprintf(formatString, time))
			continue
		}

		if e.Field(i).Type().String() == "time.Time" {
			time := e.Field(i).MethodByName("Format").Call([]reflect.Value{reflect.ValueOf(time.RFC3339Nano)})[0]
			values.WriteString(fmt.Sprintf(formatString, time))
			continue
		}

		val, err := driver.DefaultParameterConverter.ConvertValue(e.Field(i).Interface())

		if err != nil {
			return "", err
		}

		if val == nil {
			if i == e.NumField()-1 {
				values.WriteString(fmt.Sprintf("%v )", "NULL"))
				continue
			}

			values.WriteString(fmt.Sprintf("%v, ", "NULL"))
			continue
		}

		values.WriteString(fmt.Sprintf(formatString, val))

	}

	return values.String(), nil
}

func TypeConverter(val interface{}) (newval interface{}) {
	// ToDb converts val to another type. Called before INSERT/UPDATE operations
	newval = val
	return newval
}

func makeRange(min, max int) []string {
	a := make([]string, max-min+1)
	for i := range a {
		a[i] = "$" + strconv.Itoa(min+i)
	}
	return a
}
