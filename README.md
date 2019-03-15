# build-redis-data

## What

Generate redis request payload with random data of all redis commands. You can use this tool to test stability of `redis protocol` based storage engine or just generate random data for redis test.

All payload is configured by a yaml config file, (check here)[./redis_command.yml]

## Build

```
go get github.com/yongman/build-redis-data
```

## Usage

```
Usage of ./build-redis-data:
  -config string
    	redis command config file, in yaml format (default "./redis_command.yml")
  -float int
    	max float value in redis commands (default 10000000)
  -int int
    	max int value in redis commands (default 10000000)
  -key int
    	max key range number (default 100000000)
  -multi int
    	max multiple times in redis commands, such as mget/mset (default 20)
  -port int
    	redis protocol backend server, default is 6379 (default 6379)
  -server string
    	redis protocol backend server, default is 127.0.0.1 (default "127.0.0.1")
  -silent
    	run in silent mode
  -str int
    	max item or field length (default 52)
  -val int
    	max value length (default 128)
  -worker int
    	number of worker threads, default is 1 (default 1)
```

# Authors

- [tjuqxy](https://github.com/tjuqxy) Initial and most work
- [yongman](https://github.com/yongman) Arguments configurable,redis integration,parellel support.

