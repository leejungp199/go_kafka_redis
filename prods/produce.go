package prods

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type empty struct{}

type Producer struct {
	Input        chan Input
	Brokers      []string
	Topic        string
	PartitionMap map[int]empty
	close        chan bool
}

type Input struct {
	Imsi string
	Data string
}

//func (g *GProducer) Run(fg FiveGStream, topicName string) {
func (g *Producer) Run() {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("2.1.1")
	config.Version = version

	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner

	producer, err := sarama.NewAsyncProducer(g.Brokers, config)

	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Printf("Err: %v, Topic: %v, Offset: %v, Partition: %v\n", err, err.Msg.Topic, err.Msg.Offset, err.Msg.Partition)
			errors++
		}
	}()

	st := time.Now()
	dur := time.Second

ProducerLoop:
	for {
		select {
		case input := <-g.Input:
			imsiLength := len(input.Imsi)
			phoneNoLastDigit, err := strconv.Atoi(input.Imsi[imsiLength-1 : imsiLength])
			if err != nil {
				panic(err)
			}

			if _, exists := g.PartitionMap[phoneNoLastDigit]; !exists {
				// this phone number is not the one we should process. ignore it...
				continue
			}

			message := &sarama.ProducerMessage{Topic: g.Topic, Value: sarama.StringEncoder(input.Data), Partition: int32(phoneNoLastDigit)}

			producer.Input() <- message

			enqueued++
			dur := time.Now().Sub(st)

			// 15분 뒤 알아서 꺼지기
			if dur.Seconds() > 900 {
				producer.AsyncClose() // Trigger a shutdown of the producer.
				break ProducerLoop
			}

		//case <-signals:
		case <-g.close:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			fmt.Printf("Closing Producer %s...\n", g.Topic)
			time.Sleep(1 * time.Second)
			break ProducerLoop
		}
	}
	wg.Wait()
	dur = time.Now().Sub(st)
	log.Printf("Duration: %f, Enqued: %d; Successfully produced: %d; errors: %d\n", dur.Seconds(), enqueued, successes, errors)
	log.Printf("AVG: %d;\n", successes/int(dur.Seconds()))

}

func (g *Producer) Close() {
	g.close <- true
}

/*
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	g := GProducer{}
	fg := FiveGStream{
		OutputChn: make(chan []byte, 1000),
		Close:     make(chan int, 1),
	}
	g.Run(fg, "gtpc")
}
*/

func InitProducer(brokers []string, topic string, partitions []int, logAbsPath string) *Producer {
	g := Producer{
		make(chan Input, 1),
		brokers,
		topic,
		make(map[int]empty),
		make(chan bool),
	}

	// map partition index numbers that will be used to filter imsi numbers.
	// init partitionMap
	for _, i := range partitions {
		g.PartitionMap[i] = empty{}
	}

	return &g
}
