package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"github.com/redis/go-redis/v9"
)

var (
	rdbCtx = context.Background()
	c      = config.Config.RDBConfig
	rdb    *RDBConnection
)

type RDBConnection struct {
	*redis.Client
}

func ConnectRDB() (err error) {
	dbChannel, err := strconv.Atoi(c.DB)
	if err != nil {
		return err
	}
	rdb = &RDBConnection{redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Password: c.Password,
		DB:       dbChannel,
	})}
	return
}

func GetClientRDB() *RDBConnection {
	return rdb
}

func (rdb *RDBConnection) SetCache(key string, value interface{}, duration int) error {
	var timeLimit = duration
	var err error
	if duration == 0 {
		timeLimit, err = strconv.Atoi(c.Expires)
		if err != nil {
			return err
		}
	}

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	rdb.Set(rdbCtx, key, v, time.Duration(timeLimit)*time.Minute)
	return nil
}

func (rdb *RDBConnection) GetCache(key string, dest interface{}) error {
	val, err := rdb.Get(rdbCtx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (rdb *RDBConnection) DeleteCache(key string) error {
	_, err := rdb.Del(rdbCtx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *RDBConnection) GetTTLCache(key string) (int, error) {
	val, err := rdb.TTL(rdbCtx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(val.Seconds()), nil
}
