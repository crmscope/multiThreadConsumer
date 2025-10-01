#!/bin/bash

case "$1" in
  start_prod)
    cd custom/include/MultiThreadsConsumer/cmd/app && ./main
    ;;
  start)
    cd custom/include/MultiThreadsConsumer/cmd/app && go run main.go
    ;;
  stop)
    process_id=$(cat /var/run/crm/bitrix_exchange_daemon.pid)
    if [ "$process_id" = "" ]; then
        echo "The process is not running."
    else
        kill -3  "$process_id"
        echo "The call has been sent."
    fi
    ;;
  kill)
    process_id=$(cat /var/run/crm/bitrix_exchange_daemon.pid)
    if [ "$process_id" = "" ]; then
        echo "The process is not running."
    else
        kill -9  "$process_id"
        echo "The call has been sent."
    fi
    ;;
  *)
    echo "Unknown command: $cmd"
    ;;
esac
