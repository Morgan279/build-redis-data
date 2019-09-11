package data

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/yongman/build-redis-data/tools"
)

func MakeRedisData(c map[string]interface{}, redis redis.Conn, dgr *tools.DataGenerateRule) {
	// Deal conf fail
	if c == nil {
		return
	}
	dataTypes := []string{"string", "set", "zset", "hash", "pf", "list"}

	r := tools.NewRand()

	inner := func(slice []string) []interface{} {
		ret := make([]interface{}, len(slice))
		for i, arg := range slice {
			ret[i] = arg
		}
		return ret
	}

	var numKeys = 0

	// deal read command
	readCommand := c["read-command"].(map[interface{}]interface{})
	for dataType, commandConf := range readCommand {
		if dataType == "key" {
			dataType = dataTypes[r.RandInt(len(dataTypes))-1]
		}
		eachCommand := commandConf.(map[interface{}]interface{})
		for cmd, param := range eachCommand {
			paramStr := tools.BuildData(r, dgr, dataType, param, &numKeys)
			numKeys = 0
			if !dgr.Silent {
				fmt.Println(cmd, paramStr)
			}
			startTs := time.Now()
			_, err := redis.Do(cmd.(string), inner(paramStr)...)
			elapsed := time.Now().Sub(startTs)
			if dgr.Latency {
				fmt.Printf("cost: %d ms cmd: %v\n", int64(elapsed)/1000, cmd)
			}
			if err != nil {
				if !dgr.Silent {
					fmt.Println(err.Error())
				}
				if strings.Contains(err.Error(), "use of closed network connection") {
					os.Exit(1)
				}
			}
		}
	}

	// deal write command
	writeCommand := c["write-command"].(map[interface{}]interface{})
	for dataType, commandConf := range writeCommand {
		if dataType == "key" {
			dataType = dataTypes[r.RandInt(len(dataTypes))-1]
		}
		eachCommand := commandConf.(map[interface{}]interface{})
		for cmd, param := range eachCommand {
			paramStr := tools.BuildData(r, dgr, dataType, param, &numKeys)
			numKeys = 0
			if !dgr.Silent {
				fmt.Println(cmd, paramStr)
			}
			startTs := time.Now()
			_, err := redis.Do(cmd.(string), inner(paramStr)...)
			elapsed := time.Now().Sub(startTs)
			if dgr.Latency {
				fmt.Printf("cost: %d ms cmd: %v\n", int64(elapsed)/1000, cmd)
			}
			if err != nil {
				if !dgr.Silent {
					fmt.Println(err.Error())
				}
				if strings.Contains(err.Error(), "use of closed network connection") {
					os.Exit(1)
				}
			}
		}
	}

	return
}
