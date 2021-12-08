package bnats

import (
	"io/ioutil"

	"github.com/nats-io/nats.go"
)

func SendToNats(filePath string) {
	fl, err := ioutil.ReadDir(filePath)
	noErr(err)

	js := InitNats()

	for _, i := range fl {
		data, _ := ioutil.ReadFile(filePath+"/"+i.Name())
		publish(data, js)
	}
}

func publish(data []byte, js nats.JetStream) {
	_, err := js.Publish("test.test", data)
	noErr(err)
}

func InitNats() nats.JetStream {
	natsConnection, err := nats.Connect("0.0.0.0:4222")
	noErr(err)

	err = natsConnection.Flush()
	noErr(err)

	js, err := natsConnection.JetStream()
	noErr(err)

	_, err = js.AddStream(&nats.StreamConfig{
		Name: "test",
		Subjects: []string{"test.*"},
		MaxConsumers: 10,
	})

	noErr(err)

	return js
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
