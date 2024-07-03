package logger

import (
	"fmt"
	"os"
	"time"
)

func WriteIntToFile(num int) error {
	file, err := os.OpenFile("miner-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	_, err = fmt.Fprintf(file, "%s: %d\n", timestamp, num)
	return err
}
