package main

import(
	"log"
	"database/sql"

	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

var dbmap *gorp.DbMap

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", serverConfig.PgConnect)
	if err != nil {
		log.Fatalln("Failed to connect postgres:", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithNameAndSchema(Event{}, "atomic", "events").SetKeys(true, "Id")

	return dbmap
}

func postgresWriter() {
	var err error
	var event *Event

	for {
		event = <- eventsChannel
		err = dbmap.Insert()
		if err != nil {
			eventFail(counter, "Failed to write data to database", err, "", &NginxLog{}, event)
			return
		}

		successRing.Add("", &NginxLog{}, event, "success")
		log.Println(counter, "Succesfuly witten")
	}
}
