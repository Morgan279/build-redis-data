package log

import (
	"../conf"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	if *conf.LogError {
		f, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
		if err != nil {
			panic(err)
		}
		Logger = log.New(f, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	}
}
