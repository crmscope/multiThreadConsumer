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
	"encoding/json"
	"time"

	"github.com/crmscope/multiThreadConsumer/packages/kConnector"
)
import "fmt"

// This function is exported to php next jsonStr should be returned from php.
// jsonStr - is Apache topics configuration.
//
//export MultiConsume
func MultiConsume(kafkaConnect *C.char, commitInterval int16, debugMode bool, jsonStr *C.char) bool {

	var worker []kConnector.Worker
	var WStatus []kConnector.WorkerStatus

	err := json.Unmarshal([]byte(C.GoString(jsonStr)), &worker)
	if err != nil {
		return false
	}

	// Goroutine is the core of the library. There can be many threads.
	for i := 0; i < len(worker); i++ {
		WStatus = append(WStatus, kConnector.WorkerStatus{0, worker[i]})
		go kConnector.Connect(worker[i], time.Duration(commitInterval)*time.Second, &WStatus[i])
	}
	fmt.Printf("Status=%d\n", WStatus[0].Status)
	go kConnector.ManageConnect()

	// Time to complete the process
	time.Sleep(1000 * time.Second)

	return true
}

func main() {}
