package kafkaConnector

import (
	"context"
	"exchange/kafka/internal/cLoger"
	"exchange/kafka/internal/configLoader"
	"exchange/kafka/internal/fileHandler"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// Консьюмер запускается в отдельном потоке. Потоки между собой никак не связаны — каждый поток
// представляет собой самостоятельного консьюмера.
//
// Принцип работы:
//
// Консьюмер считывает данные из заданного топика и передаёт их на обработку PHP-обработчику.
// Ответ от PHP-обработчика считается успешным, если он равен 0.
// Любой другой ответ трактуется как ошибка.
// В случае ошибки указатель Offset не смещается.
func Consume(cfg *configLoader.Config, topic configLoader.Topic, wg *sync.WaitGroup, stop chan struct{}, done chan struct{}, i int) {

LOOP:
	for {

		defer wg.Done()
		// Bocker configure
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
			Topic:          topic.TopicName,
			GroupID:        topic.GroupId,
			CommitInterval: time.Duration(cfg.Kafka.CommitInterval) * time.Second, // Autocommit disable
		})
		defer r.Close()

		// Настройка логера
		loger := new(cLoger.Logger)
		loger.Init(fmt.Sprintf("%s.%d.%s", cfg.Kafka.Log, i, "log"), "[Thread] ")
		defer loger.Close()

		ctx := context.Background()

		for {
			// Читает сообщение из брокера.
			m, err := r.FetchMessage(ctx)

			if err != nil {
				// Пытается прочитать сообщение. Если в течении минуты это не удается
				// сделать возникает паника и ошибка 501 записывается в лог.
				// ВАЖНО! При отсутствии подключения консьюмер не прекращает свою работу.
				// При появлении подключения консьюмер автоматически восстановит работу.
				loger.Error(501, err.Error())
				break
			}
			// Записывает то что удалось прочитать
			loger.Info(fmt.Sprintf("brokerResponse offset=%d key=%s value=%s\n",
				m.Offset,
				string(m.Key),
				string(m.Value)))

			// Ниже идет циклически обращение к обработчику на php до тех
			// пор пока тот не вернет код 200 (Данные успешно обработаны).
			// Для обеспечения сохранности данных код упаковывается в Base64
			for {
				handlerResponse := fileHandler.FileRun(m.Value, topic.WorkerName)
				loger.Info("handlerResponse: " + string(handlerResponse))

				// Тут принцип Linux: при удачной работе программы нечего не выводится.
				// Если php-обработчик что выводит это расценивается как ошибка.
				// Для проверки удаляем все пробелы, табуляцию, символы новой строки и перевода коретки.
				if len(strings.TrimSpace(string(handlerResponse))) == 0 {
					// Autocommit
					if err := r.CommitMessages(ctx, m); err != nil {
						/**
						* Если не удается закомитить offset то отключаем поток
						* т.к. похоже с брокером проблема.
						 */
						loger.Error(503, err.Error())
						break LOOP
					}
					break
				}

				//При ошибке на стороне хендлера повторяем
				// отправку до тех пор, пока не получим код 200
				time.Sleep(2 * time.Second)
			}
		}

		select {
		case <-stop:
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

// Пока что продюсер используется только для тестовых целей.
func Produce(host string, port int, topicName string, messages []kafka.Message) {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", host, port)},
		Topic:   topicName,
		//Balancer: &kafka.RoundRobin{},
		//Balancer: &kafka.ReferenceHash{},
	})
	defer producer.Close()

	// Отправляем сообщения
	err := producer.WriteMessages(context.Background(), messages...)
	if err != nil {
		fmt.Printf("Error sending messages: %v", err)
	}
}
