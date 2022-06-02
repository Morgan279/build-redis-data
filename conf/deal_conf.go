package conf

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
)

var (
	ConfFile       = flag.String("config", "./redis_command.yml", "redis command config file, in yaml format")
	ExpectPoolHost = flag.String("es_host", "127.0.0.1", "ip of redis protocol backend server that responses expect value, default is 127.0.0.1")
	ExpectPoolPort = flag.Int("es_port", 6379, "port of redis protocol backend server that responses expect value, default is 6379")
	ExpectPoolPwd  = flag.String("es_pwd", "", "password of redis protocol backend server that responses expect value")
	ActualPoolHost = flag.String("as_host", "127.0.0.1", "ip of redis protocol backend server that responses actual value, default is 127.0.0.1")
	ActualPoolPort = flag.Int("as_port", 6379, "port of redis protocol backend server that responses actual value, default is 6379")
	ActualPoolPwd  = flag.String("as_pwd", "", "password of redis protocol backend server that responses actual value")
	Parallel       = flag.Int("worker", 1, "number of worker threads, default is 1")
	MaxKeyNum      = flag.Int("key", 100000000, "max key range number")
	MaxStrLen      = flag.Int("str", 52, "max item or field length")
	MaxValLen      = flag.Int("val", 128, "max value length")
	MaxIntVal      = flag.Int("int", 10000000, "max int value in redis commands")
	MaxFloatVal    = flag.Int("float", 10000000, "max float value in redis commands")
	MaxMulti       = flag.Int("multi", 20, "max multiple times in redis commands, such as mget/mset")
	Silent         = flag.Bool("silent", false, "run in silent mode")
	LogError       = flag.Bool("log", true, "log error information or panic directly")
	Latency        = flag.Bool("latency", false, "show latency of each command")
)

func Parse() {
	flag.Parse()
	fmt.Printf("parallel: %d\n", *Parallel)
	if *Parallel < 1 {
		fmt.Printf("parallel %d not valid, >=1 should be assigned\n", *Parallel)
		os.Exit(1)
	}
}

func LoadConfFile() map[string]interface{} {
	return dealConf(*ConfFile)
}

// read config as a map
func dealConf(confFile string) map[string]interface{} {
	c := make(map[string]interface{})

	// read file
	content, err := ioutil.ReadFile(confFile)
	if err != nil {
		fmt.Println("Read file failed:", err)
		return nil
	}

	// unmarshal
	err = yaml.Unmarshal(content, c)
	if err != nil {
		fmt.Println("Unmarshal conf failed:", err)
		return nil
	}

	return c
}
