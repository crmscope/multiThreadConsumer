package kConnector

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Performed in Goroutine.
// worker - Topic configere.
// One goroutine connected one partition in topic.
func ManageConnect() {

	// Bocker configure
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "multiThreadConsumer_Manage",
		GroupID:        "multiThreadConsumer_Manage",
		CommitInterval: 0, // Autocommit disable
	})
	defer r.Close()

	ctx := context.Background()

	for {

		// Read message
		m, err := r.FetchMessage(ctx)

		if err != nil {
			fmt.Println("error reading message:", err)
			break
		}

		fmt.Printf("offset=%d key=%s value=%s\n",
			m.Offset, string(m.Key), string(m.Value))

		// Autocommit
		if err := r.CommitMessages(ctx, m); err != nil {
			fmt.Println("error offset:", err)
		}
	}
}
