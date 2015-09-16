package main

import(
	"log"
	"database/sql"

	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

var dbmap *gorp.DbMap

func initDb() *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    pgconnect := serverConfig.Pgconnect
	db, err := sql.Open("postgres", pgconnect)
    checkErr(err, "sql.Open failed")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts

    dbmap.AddTableWithNameAndSchema(Event{}, "atomic", "events").SetKeys(true, "Id")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
