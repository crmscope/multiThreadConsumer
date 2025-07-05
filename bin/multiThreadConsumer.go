// The extension allows multiple threads to connect to Apache Kafka from php.
//
// System requirements:
// PHP ​​7.4. and higher
// PHP Extantion FFI
//
// The extension is compiled into a C module from Golang.
//
// To compile: go build -o multiThreadConsumer.so -buildmode=c-shared multiThreadConsumer.go
//
// You can find code examples in the examples folder.

package main

import (
	"C"
	"context"
	"encoding/json"
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

// This function is exported to php next jsonStr should be returned from php.
// jsonStr - is Apache topics configuration.
//
//export MultiConsume
func MultiConsume(kafkaConnect *C.char, commitInterval int16, debugMode bool, jsonStr *C.char) bool {

	var worker []Worker

	err := json.Unmarshal([]byte(C.GoString(jsonStr)), &worker)
	if err != nil {
		return false
	}

	// Goroutine is the core of the library. There can be many threads.
	for i := 0; i < len(worker); i++ {
		go KafkaConnect(worker[i])
	}

	// Time to complete the process
	time.Sleep(1000 * time.Second)

	return true
}

// Performed in Goroutine.
// worker - Topic configere.
// One goroutine connected one partition in topic.
func KafkaConnect(worker Worker) {

	script := "worker.php"

	// Bocker configure
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          worker.Name,
		GroupID:        "my-group",
		CommitInterval: 0, // Autocommit disable
	})
	defer r.Close()

	ctx := context.Background()

	for {

		//m, err := r.ReadMessage(ctx)
		m, err := r.FetchMessage(ctx)

		if err != nil {
			fmt.Println("error reading message:", err)
			break
		}

		fmt.Printf("offset=%d key=%s value=%s\n",
			m.Offset, string(m.Key), string(m.Value))

		args := []string{script, "REZZ:", string(m.Value)}
		cmd := exec.Command("php", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка запуска PHP: %v\n", err)
			return
		}
		// Print result
		fmt.Printf("Print: %s\n", out)

		// Autocommit
		if err := r.CommitMessages(ctx, m); err != nil {
			fmt.Println("error offset:", err)
		}

	}

	fmt.Println(worker)
}

func main() {}
