package fileHandler

// Пакет передает полученные данные на стандартный вход PHP-обработчика.
// Все что получено от обработчика возвращается как есть в точку вызова.
// Для сохранения целостности данных, они упаковываются в Base64.

import (
	"encoding/base64"
	"os/exec"
)

// Путь к папке с воркерами
const WORKER_PATH = "./../../handlers/"

func FileRun(message []byte, WorkerName string) []byte {
	cmd := exec.Command("php", WORKER_PATH+WorkerName)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	// Упаковывает данные в Base64
	encoded := base64.StdEncoding.EncodeToString(message)

	// Пишем данные в stdin
	stdin.Write([]byte(encoded))
	stdin.Close()

	// Читаем результат из stdout
	buf := make([]byte, 1024)
	n, _ := stdout.Read(buf)

	cmd.Wait()
	return buf[:n]
}
