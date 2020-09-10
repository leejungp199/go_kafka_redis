package module

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	verbose = false
	handler sarama.ConsumerGroupHandler
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Brokers []string
	Version string
	Topic   string
	Cons    string
}

// Kafka Session handler function for
func (c *Consumer) RunConsGroup() {
	/// prepare consumer
	log.Println("Starting a new Sarama consumer")

	//define brokers and topic
	brokers := c.Brokers
	topic := c.Topic

	version, err := sarama.ParseKafkaVersion(c.Version)

	if err != nil {
		panic(err)
	}
	//get current kafka version

	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = version

	//offset 어디서 부터 읽을지 설정
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	consGroup := c.Cons

	group, err := sarama.NewConsumerGroup(brokers, consGroup, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	//입력 모듈에 해당하는 consumer group handler 실행
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		if consGroup == "session" {
			handler = sessionGroupHandler{}
		} else if consGroup == "ex" {
			handler = exGrouphandler{}
		} else {
			fmt.Println("ERROR no matched group name", err)
		}

		//kafka 서버로 부터 데이터를 수신해 group handler로 보낸다
		err := group.Consume(ctx, strings.Split(topic, ","), handler)
		if err != nil {
			panic(err)
		}
	}
}

func InitConsumer(brokers []string, topic string, cons string) *Consumer {
	//logName := fmt.Sprintf("CONSUMER_%s_%d", topic, partition)
	//logPath := fmt.Sprintf("%s/consumer_%s_%d.log", logAbsPath, topic, partition)

	consumer := Consumer{
		Brokers: brokers,
		Version: "2.1.1",
		Topic:   topic,
		Cons:    cons,
	}
	return &consumer
}
