package tools

import (
	"../conf"
	"../log"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"sync"
)

var once sync.Once
var expectPool *redis.Pool
var actualPool *redis.Pool

func getInstance() (*redis.Pool, *redis.Pool) {
	once.Do(func() {
		fmt.Printf("connecting to expect value pool %s:%d\n", *conf.ExpectPoolHost, *conf.ExpectPoolPort)
		expectPool = NewPool(fmt.Sprintf("%s:%d", *conf.ExpectPoolHost, *conf.ExpectPoolPort), *conf.ExpectPoolPwd)
		fmt.Printf("connecting to actual value pool %s:%d\n", *conf.ActualPoolHost, *conf.ActualPoolPort)
		actualPool = NewPool(fmt.Sprintf("%s:%d", *conf.ActualPoolHost, *conf.ActualPoolPort), *conf.ActualPoolPwd)
	})
	return expectPool, actualPool
}

func Judge(commandName string, args []interface{}) error {
	expectPool, actualPool = getInstance()

	epRes, epErr := expectPool.Get().Do(commandName, args...)

	apRes, apErr := actualPool.Get().Do(commandName, args...)

	if epErr != nil {
		if !reflect.DeepEqual(epErr, apErr) {
			errStr := fmt.Sprintf("inconsistent err response:\n"+
				"expect err response: %+v\n"+
				"but actual err response is: %+v\n"+
				"executed command: %s, args: %+v\n",
				epErr, apErr, commandName, args)
			handleErr(errStr)
		}
		return epErr
	}

	if apErr != nil {
		errStr := fmt.Sprintf("unexpected err response:\n"+
			"expect response: %+v\n"+
			"but actuall responses an error: %+v\n"+
			"executed command: %s, args: %+v\n",
			epRes, apErr, commandName, args)
		handleErr(errStr)
		return apErr
	}

	if !reflect.DeepEqual(apRes, epRes) {
		errStr := fmt.Sprintf("unexpected response:\n"+
			"actual: %+v\n"+
			"expect: %+v\n"+
			"executed command: %s, args: %+v\n",
			apRes, epRes, commandName, args)

		handleErr(errStr)
	}

	return nil
}

func handleErr(errStr string) {
	if *conf.LogError {
		fmt.Println(errStr)
		log.Logger.Println(errStr)
	} else {
		panic(errStr)
	}
}
