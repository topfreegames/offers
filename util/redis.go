// Package util for redis helpers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>
package util

import (
	"fmt"

	redis "gopkg.in/redis.v5"

	"github.com/Sirupsen/logrus"
)

// RedisClient identifies uniquely one redis client with a pool of connections
type RedisClient struct {
	Logger logrus.FieldLogger
	Client *redis.Client
}

// GetRedisClient creates and returns a new redis client based on the given settings
func GetRedisClient(redisHost string, redisPort int, redisPassword string, redisDB int, maxPoolSize int, logger logrus.FieldLogger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
		PoolSize: maxPoolSize,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	cl := &RedisClient{
		Client: client,
		Logger: logger,
	}
	return cl, nil
}

// GetConnection return a redis connection
func (c *RedisClient) GetConnection() *redis.Client {
	return c.Client
}
