package smongo

import (
	"context"
	"encoding/json"
	"fmt"
	"renamosy/bnats"
	"renamosy/bpanda"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	PeriodStartUtc string `json:"period_start_utc"`
	PeriodEndUtc   string `json:"period_end_utc"`
	ReadingType    string `json:"reading_type"`
	Value          int    `json:"value"`
}

func FetchFromNats() {
	js := bnats.InitNats()

	js.Subscribe("test.test", mongoHandler)
}

func FetchFromRedpanda() {
	r, _ := bpanda.InitRedpanda()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var msg message

		err = json.Unmarshal(m.Value, &msg)
		noErr(err)

		coll := initMongo()

		input := make([]interface{}, len(msg.Timeseries))
		for i := range msg.Timeseries {
			input[i] = msg.Timeseries[i]
		}

		t1 := time.Now()
		_, err = coll.InsertMany(context.Background(), input)
		noErr(err)

		fmt.Printf("%+v Milliseconds to insert %d ts from redpanda to mongo\n", time.Now().Sub(t1).Milliseconds(), len(input))
	}
}

func mongoHandler(msg *nats.Msg) {
	coll := initMongo()

	var m message
	err := json.Unmarshal(msg.Data, &m)
	noErr(err)

	input := make([]interface{}, len(m.Timeseries))
	for i := range m.Timeseries {
		input[i] = m.Timeseries[i]
	}

	t1 := time.Now()
	_, err = coll.InsertMany(context.Background(), input)
	noErr(err)

	fmt.Printf("%+v Milliseconds to insert %d ts from nats to mongo\n", time.Now().Sub(t1).Milliseconds(), len(input))

	err = msg.Ack()
	noErr(err)
}

func initMongo() *mongo.Collection {
	c, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	noErr(err)

	err = c.Connect(context.Background())
	noErr(err)

	testColl := c.Database("SHELL").Collection("test")
	return testColl
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
