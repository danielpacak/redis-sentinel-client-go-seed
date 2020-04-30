package redis

import (
	"errors"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/etc"
	xredis "github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func GetPool(config etc.Config) (pool *xredis.Pool) {
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{config.RedisAddr},
		MasterName: config.RedisMaster,
		Dial: func(addr string) (conn xredis.Conn, err error) {
			log.WithField("sentinel_addr", config.RedisAddr).Debug("Connecting to Redis sentinel")
			conn, err = xredis.Dial("tcp", addr,
				xredis.DialConnectTimeout(config.RedisConnectTimeout),
				xredis.DialReadTimeout(config.RedisReadTimeout),
				xredis.DialWriteTimeout(config.RedisWriteTimeout))
			if err != nil {
				log.WithError(err).Error("Error while connecting to Redis sentinel")
				return
			}
			return
		},
	}

	pool = &xredis.Pool{
		MaxIdle:     config.RedisPoolMaxIdle,
		MaxActive:   config.RedisPoolMaxActive,
		IdleTimeout: config.RedisPoolIdleTimeout,

		Dial: func() (conn xredis.Conn, err error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return
			}
			log.WithField("master_addr", masterAddr).Debug("Connecting to Redis master")
			conn, err = xredis.Dial("tcp", masterAddr)
			if err != nil {
				log.WithError(err).Error("Error while connecting to Redis master")
				return
			}
			return
		},

		TestOnBorrow: func(c xredis.Conn, t time.Time) (err error) {
			log.Debugf("Testing connection on borrow: %v", t)
			if !sentinel.TestRole(c, "master") {
				err = errors.New("role check failed")
				log.WithError(err).Error("Error while testing connection on borrow")
				return
			}
			return
		},
	}

	pool.Stats()

	return
}
