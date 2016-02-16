package main

import(
	"log"
	"database/sql"

	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

type Postgres struct {
	ConnChan     chan *gorp.DbMap
	ConnFailChan chan bool
}

func (postgres *Postgres) Connect(config *Worker.Config) {
	db, err := sql.Open("postgres", config.PgConnect)
	if err != nil {
		postgres.ConnFailChan <- true
		return
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithNameAndSchema(SnowplowEvent{}, "atomic", "events").SetKeys(true, "Id")
	dbmap.AddTableWithNameAndSchema(AdjustEvent{}, "adjust", "events").SetKeys(true, "Id")

	postgres.ConnChan <- dbmap
}

func checkStringForNull(eventStr string, event *sql.NullString) {
	if len(eventStr) == 0 {
		event.Valid = false
	} else {
		event.String = eventStr
		event.Valid = true
	}
}

func checkIntForNull(eventInt string, event *sql.NullInt64) {
	var err error

	if eventInt == "" {
		event.Valid = false
	} else {
		event.Int64, err = strconv.ParseInt(eventInt, 10, 64)
		if err != nil {
			event.Valid = false
			return
		}

		event.Valid = true
	}
}

func checkFloatForNull(eventFloat string, event *sql.NullFloat64) {
	var err error

	if eventFloat == "" {
		event.Valid = false
	} else {
		event.Float64, err = strconv.ParseFloat(eventFloat, 64)
		if err != nil {
			event.Valid = false
			return
		}

		event.Valid = true
	}
}
