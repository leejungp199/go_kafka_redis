package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"./module"
	"./module/config"
	"./prods"
)

var (
	w   sync.WaitGroup
	cfg *config.Config
)

func main() {
	// set max MAX PROCS to number of cpu
	//멀티 프로세싱을 위한 cpu 개수 설정
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 3 {
		fmt.Println("Error. Missing configuration path argument.")
		fmt.Println("Ex) go run main.go {path to configuration file} {prod/cons}")
		os.Exit(1)
	}

	//// set config file path
	configPath := os.Args[1]
	op := os.Args[2]

	// init configs
	cfg = initConfig(configPath)

	if op == "sess" {
		// prepare consumer
		brokers := cfg.Kafka.Brokers
		topic := strings.ToLower(fmt.Sprintf("c")) // gtpC

		//consumer type: session
		cons := "sess"

		//initialize and run consumer
		consumer := module.InitConsumer(brokers, topic, cons)
		consumer.RunConsGroup()

	} else if op == "si" {
		// prepare consumer
		brokers := cfg.Kafka.Brokers
		topic := strings.ToLower(fmt.Sprintf("b")) // gtpB

		//consumer type: si
		cons := "si"

		//initialize and run consumer
		consumer := module.InitConsumer(brokers, topic, cons)
		consumer.RunConsGroup()

	} else if op == "prodb" {
		topicBB := "bearer"
		// start producer
		runProd(topicBB)

	} else {
		panic("Illegal Argument =  " + op)
	}
	w.Wait()
}

//// SIGTERM(ctrl-c) 보내면 모든 프로세스 종료
func handleSignal(handler func()) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals
	go func() {
		select {
		case <-sigterm:
			handler()
		}
	}()
}

// Kafka Producer handler function
func runProd(isGTP string) {

	if isGTP == "c" {
		w.Add(1)

		// prepare C Topic name as  NAMEc
		topic := strings.ToLower(fmt.Sprintf("c")) /
		// initialize kafka producer
		prod := prods.InitProducer(cfg.Kafka.Brokers, topic, cfg.Kafka.Partitions, cfg.Kafka.Log_folder_name)

		// run producer in a seperate goroutine
		go func() {
			defer w.Done()
			prod.Run()
		}()

		// start a goroutine that generates sample GTP-C data
		go func() {
			//module.GenerateControl(cfg.StreamData.DataCCC, prod.Input, cfg.GTPCFieldIndex.IdxImsi, false)
			module.GeneratePath(cfg.StreamData.DataCCC, prod.Input, cfg.CFieldIndex.IdxImsi, false)
		}()

		// Handle interrupt
		handleSignal(func() {
			// Order is important! Session must be closed first. Then consumer...
			prod.Close()
		})
	} else if isGTP == "b" {
		w.Add(1)

		// prepare GTP-U Topic name as  NAMEu
		topic := strings.ToLower(fmt.Sprintf("b")) // b
		// initialize kafka producer
		prod := prods.InitProducer(cfg.Kafka.Brokers, topic, cfg.Kafka.Partitions, cfg.Kafka.Log_folder_name)

		// run producer in a seperate goroutine
		go func() {
			defer w.Done()
			prod.Run()
		}()

		// start a goroutine that generates sample GTP-U data
		go func() {
			module.Generate(cfg.StreamData.DataB, prod.Input, cfg.BFieldIndex.IdxImsi, false)
		}()

		// Handle interrupt
		handleSignal(func() {
			// Order is important! Session must be closed first. Then consumer...
			prod.Close()
		})
	} else {
		w.Add(1)
		// prepare U Topic name as  NAMEu
		topic := strings.ToLower(fmt.Sprintf("t")) // 

		// initialize kafka producer
		prod := prods.InitProducer(cfg.Kafka.Brokers, topic, cfg.Kafka.Partitions, cfg.Kafka.Log_folder_name)

		// run producer in a seperate goroutine
		go func() {
			defer w.Done()
			prod.Run()
		}()

		// start a goroutine that generates sample GTP-U data
		go func() {
			module.Generate(cfg.StreamData.DataUT, prod.Input, cfg.UFieldIndex.IdxUserIp, false)
		}()

		// Handle interrupt
		handleSignal(func() {
			// Order is important! Session must be closed first. Then consumer...
			prod.Close()
		})
	}
}
