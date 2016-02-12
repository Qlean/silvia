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

var dbmap *gorp.DbMap

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
