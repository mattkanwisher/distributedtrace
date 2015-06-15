package zipkin

import (
	"encoding/json"

	"gopkg.in/redis.v3"
)

type redisOutput struct {
	*redis.Client
	config *Config
	key    string
}

// NewRedisOutput() returns an Output that converts ZipKin spans to JSON and RPUSH
// (http://redis.io/commands/RPUSH) it to the specified Redis (http://redis.io) server.
func NewRedisOutput(config *Config, addr, password string, db int64, key string) (Output, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if e := client.Ping().Err(); e != nil {
		return nil, e
	}

	return &redisOutput{client, config, key}, nil
}

func (ro *redisOutput) Write(result OutputMap) (e error) {
	var buffer []byte
	if buffer, e = json.Marshal(result); e != nil {
		return e
	}

	if _, e = ro.RPush(ro.key, string(buffer)).Result(); e != nil {
		return e
	}

	return nil
}
