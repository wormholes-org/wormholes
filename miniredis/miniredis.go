package miniredis

import (
	"encoding/json"
	"regexp"
	"github.com/go-redis/redis/v7"
)

// Ethereum interface
type Ethereum interface {
	PrintAll() string
}

var (
	eth Ethereum
	addr      = "127.0.0.1:6379"
	client    *redis.Client
	client1    *redis.Client
	logCh     = make(chan map[string]interface{}, 1000)
	nodiscard = false
	// save ip
)

func runminiredisNodiscardLog() {
	client = redis.NewClient(&redis.Options{Addr: addr, Password: "", DB: 0})
	client1 = redis.NewClient(&redis.Options{Addr: addr, Password: "", DB: 1})
	reg, _ := regexp.Compile("^\\d+$")
	for {
		select {
		case logMap := <-logCh:
			for k, v := range logMap {
				if reg.MatchString(k) {
					data, _ := json.Marshal(v)
					client.SAdd(k, data)
				} else {
					client1.SAdd(k, v)
				}
			}
		}
	}
}

func runminiredisDiscardLog() {
	client = redis.NewClient(&redis.Options{Addr: addr, Password: "", DB: 0})
	for {
		select {
		case <-logCh:
		}
	}
}

// Newminiredis Create redis client
func Newminiredis(nodiscard bool) {
	if nodiscard {
		go runminiredisNodiscardLog()
	} else {
		go runminiredisDiscardLog()
	}
}

// GetData from redis
func GetData(pKey string) (interface{}, error) {
	result, err := client.Do("GET", pKey).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}

// SetData to redis
func SetData(pKey string, pValue string) {
	client.Do("SET", pKey, pValue)
}

// lpush data to redis
func LpushData(pkey string, pVlaue string) {
	client.LPush(pkey, pVlaue)
}

// rpush data to redis
func RpushData(pkey string, pVlaue string) {
	client.RPush(pkey, pVlaue)
}
func SAdd(pkey string, pValue string) {
	client.SAdd(pkey, pValue)
}

func GetLogCh() chan map[string]interface{} {
	return logCh
}
