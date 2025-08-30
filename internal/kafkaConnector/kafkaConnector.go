package kafkaConnector

import (
	"context"
	"encoding/json"
	"exchange/kafka/internal/cLoger"
	"exchange/kafka/internal/configLoader"
	"exchange/kafka/internal/fileHandler"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

/**
 * Consume reads messages from a Kafka topic and processes them using a PHP handler.
 * It commits the message if the handler returns a successful response.
 * If the handler returns an error, it retries after a delay.
 */
type workerResponse struct {
	Error        int    `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func Consume(cfg *configLoader.Config, topic configLoader.Topic, wg *sync.WaitGroup) {

	defer wg.Done()
	// Bocker configure
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
		Topic:          topic.TopicName,
		GroupID:        topic.GroupId,
		CommitInterval: time.Duration(cfg.Kafka.CommitInterval) * time.Second, // Autocommit disable
		// SessionTimeout:    100 * time.Second,
		// HeartbeatInterval: 30 * time.Second,
		// MaxWait:           100 * time.Second,
		// RebalanceTimeout:  120 * time.Second,
		// MaxAttempts:       10,
		// StartOffset:       kafka.FirstOffset,
	})
	defer r.Close()

	// Loger configure
	logr, err := cLoger.NewFileLogger(cfg.Kafka.Log, "[MultiThreads] ")
	if err != nil {
		log.Fatalf("Logger initialization error: %v", err)
	}
	defer logr.Close()

	ctx := context.Background()

	for {
		// Read message
		m, err := r.FetchMessage(ctx)

		if err != nil {
			logr.Error(fmt.Sprint("Error reading message:", err))
			break
		}

		fmt.Printf("offset=%d key=%s value=%s\n",
			m.Offset, string(m.Key), string(m.Value))

		for {
			dataFromBroker := fileHandler.FileRun(m.Value)

			var response workerResponse
			if err = json.Unmarshal(dataFromBroker, &response); err != nil {
				logr.Error(fmt.Sprintln("JSON invalid:", err))
				return
			}

			fmt.Println("--------------PHP RESPONSE----------------")
			fmt.Println(response)

			if response.Error == 200 {
				// Autocommit
				if err := r.CommitMessages(ctx, m); err != nil {
					logr.Error(fmt.Sprintln("Error offset:", err))
					return
				}
				break
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func Produce(host string, port int, topicName string, messages []kafka.Message) {

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", host, port)},
		Topic:   topicName,
		//Balancer: &kafka.RoundRobin{},
		//Balancer: &kafka.ReferenceHash{},
	})
	defer producer.Close()

	err := producer.WriteMessages(context.Background(), messages...)
	if err != nil {
		log.Fatalf("Error sending messages: %v", err)
	}

}
