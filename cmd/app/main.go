package main

import (
	"exchange/kafka/internal/cLoger"
	"exchange/kafka/internal/configLoader"
	"exchange/kafka/internal/kafkaConnector"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

var (
	stop = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
)

const (
	CONFIG_OVERRIDE_FILE_PATH = "./../../../../../"
	CONFIG_OVERRIDE_FILE_NAME = "kafka_config.yaml"
)

func main() {

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	// Загружает yaml-файл конфига
	cfg, err := configLoader.LoadConfig(CONFIG_OVERRIDE_FILE_PATH + CONFIG_OVERRIDE_FILE_NAME)
	if err != nil {
		fmt.Println("Config file error: ", err)
		return
	}

	// Настройка логера
	loger := new(cLoger.Logger)
	loger.Init(cfg.Kafka.DaemonLogFile, "[Daemon] ")
	defer loger.Close()

	var (
		signal = flag.String("s", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown
  reload — reloading the configuration file`)
	)

	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: cfg.Kafka.DaemonPidFile,
		PidFilePerm: 0644,
		LogFileName: cfg.Kafka.DaemonLogFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{cfg.Kafka.ProcessName},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			loger.Error(510, err.Error())

		}
		err = daemon.SendCommands(d)
		if err != nil {
			loger.Error(511, err.Error())
		}
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		loger.Error(512, err.Error())
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	loger.Info("daemon started")

	wg := new(sync.WaitGroup)
	i := 0

	// Start consumers for each topic
	for _, topic := range cfg.Kafka.Topics {
		wg.Add(1)
		go kafkaConnector.Consume(cfg, topic, wg, stop, done, i)
		i++
	}
	wg.Wait()

	err = daemon.ServeSignals()
	if err != nil {
		loger.Error(513, err.Error())
	}

	loger.Info("daemon terminated")
}

func termHandler(sig os.Signal) error {
	select {
	case stop <- struct{}{}:
		// успешно отправили сигнал
	default:
		// сигнал уже был отправлен
	}

	// При сигнале kill 3  - корректно заврешаем работу
	if sig == syscall.SIGQUIT {
		<-done
	}

	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	return nil
}
