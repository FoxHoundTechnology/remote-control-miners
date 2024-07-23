package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var fileMutex sync.Mutex

func WriteIntToFile(num int, areaCode uint) error {
	// Lock the mutex before entering the critical section
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile("miner-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	_, err = fmt.Fprintf(file, "%s: Area Code - %d Container - %d\n", timestamp, areaCode, num)
	return err
}
