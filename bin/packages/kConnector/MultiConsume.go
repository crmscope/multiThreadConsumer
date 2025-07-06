package kConnector

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/segmentio/kafka-go"
)

// Structure for json configuration from php.
type Worker struct {
	Name      string `json:"name"`
	Partition int    `json:"partition"`
	Worker    string `json:"worker"`
	GroupId   string `json:"group_id"`
}

type WorkerStatus struct {
	Status int
	Worker Worker
}

// Performed in Goroutine.
// worker - Topic configere.
// One goroutine connected one partition in topic.
func Connect(worker Worker, commitInterval time.Duration, WStatus *WorkerStatus) {

	// Bocker configure
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          worker.Name,
		GroupID:        worker.GroupId,
		CommitInterval: commitInterval, // Autocommit disable
	})
	defer r.Close()

	ctx := context.Background()

	for {

		// Read message
		m, err := r.FetchMessage(ctx)
		fmt.Printf("Status=%d\n", WStatus.Status)
		if err != nil {
			fmt.Println("error reading message:", err)
			WStatus.Status = 0
			break
		}

		fmt.Printf("offset=%d key=%s value=%s\n",
			m.Offset, string(m.Key), string(m.Value))

		args := []string{worker.Worker, "REZZ:", string(m.Value)}
		cmd := exec.Command("php", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка запуска PHP: %v\n", err)
			WStatus.Status = 0
			return
		}
		// Print result
		fmt.Printf("Print: %s\n", out)

		// Autocommit
		if err := r.CommitMessages(ctx, m); err != nil {
			fmt.Println("error offset:", err)
		}

		WStatus.Status = 1

	}

	fmt.Println(worker)
}
