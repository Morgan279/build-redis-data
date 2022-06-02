package main

import (
	"./conf"
	"./data"
	"./tools"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func main() {
	conf.Parse()

	dgr := &tools.DataGenerateRule{
		MaxKeyNum:   *conf.MaxKeyNum,
		MaxStrLen:   *conf.MaxStrLen,
		MaxValLen:   *conf.MaxValLen,
		MaxIntVal:   *conf.MaxIntVal,
		MaxFloatVal: *conf.MaxFloatVal,
		MaxMulti:    *conf.MaxMulti,
		Silent:      *conf.Silent,
		Latency:     *conf.Latency,
	}

	var wg sync.WaitGroup
	sc := make(chan os.Signal, 1)
	quit := false

	signal.Notify(sc, os.Interrupt)
	go func() {
		sig := <-sc
		fmt.Printf("Got signal [%s] to exit\n", sig)
		quit = true
		os.Exit(0)
	}()

	c := conf.LoadConfFile()
	for i := 0; i < *conf.Parallel; i++ {
		wg.Add(1)
		go func() {
			for {
				if quit {
					wg.Done()
					return
				}
				data.MakeRedisData(c, dgr)
			}
		}()
	}
	wg.Wait()
	return
}
