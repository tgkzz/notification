package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"
	"github.com/tgkzz/notification/dto"
	kafkaApp "github.com/tgkzz/notification/internal/app/kafka"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:51807",
		"acks":              "all",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	key := "whatsapp"

	topic := kafkaApp.Topic

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	for n := 0; n < 20; n++ {
		data := dto.WhatsappPayload{To: "+77765922520", Message: faker.Sentence()}
		req, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          req,
		}, nil); err != nil {
			fmt.Println(err)
			return
		}
	}

	p.Flush(15 * 1000)
	p.Close()
}
