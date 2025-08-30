package main

import (
	"exchange/kafka/internal/configLoader"
	"exchange/kafka/internal/kafkaConnector"
	"fmt"
	"sync"
)

func main() {
	const configFilePath = "./../../../../../"
	const configFileName = "config_kafka.yaml"

	// Laod configuration
	cfg, err := configLoader.LoadConfig(configFilePath + configFileName)
	if err != nil {
		fmt.Println("Config file error: ", err)
		return
	}

	wg := new(sync.WaitGroup)
	i := 0
	// Start consumers for each topic
	for _, topic := range cfg.Kafka.Topics {
		wg.Add(1)
		go kafkaConnector.Consume(cfg, topic, wg)
		i++
	}
	wg.Wait()

}
