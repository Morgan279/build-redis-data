package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/yongman/build-redis-data/conf"
	"github.com/yongman/build-redis-data/data"
	"github.com/yongman/build-redis-data/tools"
)

var (
	confFile    = flag.String("config", "./redis_command.yml", "redis command config file, in yaml format")
	ip          = flag.String("server", "127.0.0.1", "redis protocol backend server, default is 127.0.0.1")
	port        = flag.Int("port", 6379, "redis protocol backend server, default is 6379")
	parallel    = flag.Int("worker", 1, "number of worker threads, default is 1")
	maxKeyNum   = flag.Int("key", 100000000, "max key range number")
	maxStrLen   = flag.Int("str", 52, "max item or field length")
	maxValLen   = flag.Int("val", 128, "max value length")
	maxIntVal   = flag.Int("int", 10000000, "max int value in redis commands")
	maxFloatVal = flag.Int("float", 10000000, "max float value in redis commands")
	maxMulti    = flag.Int("multi", 20, "max multiple times in redis commands, such as mget/mset")
	silent      = flag.Bool("silent", false, "run in silent mode")
	password    = flag.String("passwd", "", "redis auth password")
)

func main() {
	flag.Parse()

	dgr := &tools.DataGenerateRule{
		MaxKeyNum:   *maxKeyNum,
		MaxStrLen:   *maxStrLen,
		MaxValLen:   *maxValLen,
		MaxIntVal:   *maxIntVal,
		MaxFloatVal: *maxFloatVal,
		MaxMulti:    *maxMulti,
		Silent:      *silent,
	}

	var wg sync.WaitGroup
	sc := make(chan os.Signal, 1)
	quit := false
	fmt.Println(*ip, *port, *parallel)

	if *parallel < 1 {
		fmt.Println("parallel not valid, >=1 should be assigned")
		os.Exit(1)
	}
	signal.Notify(sc, os.Interrupt)
	go func() {
		sig := <-sc
		fmt.Println("Got signal [%s] to exit", sig)
		quit = true
		os.Exit(0)
	}()

	pool := tools.NewPool(fmt.Sprintf("%s:%d", *ip, *port), *password)

	c := conf.DealConf(*confFile)
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			redis := pool.Get()
			for {
				if quit {
					wg.Done()
					return
				}
				data.MakeRedisData(c, redis, dgr)
			}
		}()
	}
	wg.Wait()
	return
}
