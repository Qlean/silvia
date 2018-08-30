package silvia

import (
	"database/sql"
	"strconv"

	"gopkg.in/gorp.v1"
)

type Postgres struct {
	Connection *gorp.DbMap
}

type Redshift struct {
	Connection   *gorp.DbMap
	snowplowTmap *gorp.TableMap
	adjustTmap   *gorp.TableMap
}

func (redshift *Redshift) Connect(config *Config) error {

	db, err := sql.Open("postgres", config.RedshiftConnect)
	if err != nil {
		return err
	}

	// snowplowTmap.Columns.Columnname
	// db.Prepare(pq.CopyInSchema("snowplow", "events", GetColumns(snowplowTmap)...))

	// stmt, err := db.Prepare(pq.CopyIn("users", "name", "age"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// snowplowTmap := dbmap.AddTableWithNameAndSchema(SnowplowEvent{}, "atomic", "events")
	// adjustTmap := dbmap.AddTableWithNameAndSchema(AdjustEvent{}, "adjust", "events")

	redshift.Connection = dbmap
	return nil
}

// GetColumns fetches column names from gorp tablemap
func GetColumns(tmap *gorp.TableMap) (columns []string) {

	for index, _ := range tmap.Columns {
		columns = append(columns, tmap.Columns[index].ColumnName)
	}
	return columns
}

func (postgres *Postgres) Connect(config *Config) error {
	db, err := sql.Open("postgres", config.PostgresConnect)
	if err != nil {
		return err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithNameAndSchema(SnowplowEvent{}, "atomic", "events").SetKeys(true, "Id")
	dbmap.AddTableWithNameAndSchema(AdjustEvent{}, "adjust", "events").SetKeys(true, "Id")

	postgres.Connection = dbmap
	return nil
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

func (event *SnowplowEvent) Write(dbmap *gorp.DbMap) error {
	err := dbmap.Insert(event)
	if err != nil {
		return err
	}

	return nil
}
