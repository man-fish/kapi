package utils

import (
	"fmt"
	"os"
	"time"
)

func todayLogFile() string {
	return fmt.Sprintf("%v/log/log_%v.txt",RootPath(),time.Now().Format("2006-01-01 15:04:05"))
}

func NewLogFile() (f *os.File, err error) {
	return os.OpenFile(todayLogFile(),os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)
}