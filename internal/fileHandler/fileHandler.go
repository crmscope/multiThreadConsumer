package fileHandler

import (
	"fmt"
	"os/exec"
)

func FileRun(message []byte) []byte {
	cmd := exec.Command("php", "./../../handlers/HandlerOpportunity.php")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	stdin.Write(message)
	stdin.Close()

	buf := make([]byte, 1024)
	n, _ := stdout.Read(buf)
	fmt.Println(string(buf[:n]))

	cmd.Wait()
	return buf[:n]
}
