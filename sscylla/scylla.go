package sscylla

import (
	"encoding/json"
	"fmt"
	"renamosy/bnats"
	"time"

	"github.com/gocql/gocql"
	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

type message struct {
	CountryCodeAlpha3  string `json:"country_code_alpha3"`
	MeterID            string `json:"meter_id"`
	DataProvider       string `json:"data_provider"`
	SourceRef          string `json:"source_ref"`
	SourceTimestampUtc string `json:"source_timestamp_utc"`
	Granularity        string `json:"granularity"`
	Direction          string `json:"direction"`
	Units              string `json:"units"`
	Timeseries         []ts   `json:"timeseries"`
}

type ts struct {
	PeriodStartUtc string  `json:"period_start_utc"`
	PeriodEndUtc   string  `json:"period_end_utc"`
	ReadingType    string  `json:"reading_type"`
	Value          float64 `json:"value"`
}

var initScylladb gocqlx.Session

func FetchFromNats() {
	initScylladb = InitScylla()

	js := bnats.InitNats()
	js.Subscribe("test.test", scyllaHandler)
}

func scyllaHandler(msg *nats.Msg) {
	var m message

	err := json.Unmarshal(msg.Data, &m)
	noErr(err)

	var tsMetadata = table.Metadata{
		Name:    "shell.ts",
		Columns: []string{"period_start_utc", "period_end_utc", "reading_type", "value"},
		PartKey: []string{"period_start_utc"},
		SortKey: []string{"period_start_utc"},
	}

	var personTable = table.New(tsMetadata)

	t1 := time.Now()

	for _, tsss := range m.Timeseries {
		q := initScylladb.Query(personTable.Insert()).BindStruct(tsss)
		if err := q.ExecRelease(); err != nil {
			panic(err)
		}
	}

	fmt.Printf("%+v Milliseconds to insert %d ts from nats to scylla\n", time.Now().Sub(t1).Milliseconds(), len(m.Timeseries))

	err = msg.Ack()
	noErr(err)
}

func InitScylla() gocqlx.Session {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	// cluster.Keyspace = "shell"

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	noErr(err)

	session.ExecStmt("CREATE KEYSPACE shell WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy'};")

	err = session.ExecStmt(`CREATE TABLE shell.ts (
   period_start_utc text,
   period_end_utc text,
   reading_type text,
   Value text,
   PRIMARY KEY(period_start_utc));`)
	if err != nil {
		println(err.Error())
	}
	// noErr(err)

	return session
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
