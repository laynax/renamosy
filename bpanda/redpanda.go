package bpanda

import (
	"context"
	"io/ioutil"

	"github.com/segmentio/kafka-go"
)

func SendToRedpanda(filePath string) {
	_, writer := InitRedpanda()

	fl, err := ioutil.ReadDir(filePath)
	noErr(err)

	for _, i := range fl {
		data, _ := ioutil.ReadFile(filePath + "/" + i.Name())

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("topic"),
				Value: data,
			},
		)
		noErr(err)
	}
}

func InitRedpanda() (*kafka.Reader, *kafka.Writer) {
	url := "0.0.0.0:9092"
	_, err := kafka.Dial("tcp", url)
	noErr(err)

	w := &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    "test.test",
		Balancer: &kafka.LeastBytes{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{url},
		Topic:     "test.test",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return r, w
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
