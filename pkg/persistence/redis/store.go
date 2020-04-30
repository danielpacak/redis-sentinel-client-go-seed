package redis

import (
	"strings"

	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/persistence"
	xredis "github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type store struct {
	pool *xredis.Pool
}

func NewStore(pool *xredis.Pool) persistence.Store {
	return &store{
		pool: pool,
	}
}

func (s *store) Set(key, value string) (err error) {
	conn := s.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("SET", key, value)
	if err != nil {
		return
	}
	log.Tracef("Reply: %v", reply)
	return
}

func (s *store) Keys() (keys []string, err error) {
	conn := s.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("KEYS", "*")
	if err != nil {
		return
	}
	log.Tracef("Reply: %v", reply)
	keys, err = xredis.Strings(reply, err)
	return
}

func (s *store) Info() (infos []string, err error) {
	conn := s.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("INFO", "all")
	if err != nil {
		return
	}
	log.Tracef("Reply: %v", reply)
	info, err := xredis.String(reply, err)
	if err != nil {
		return
	}
	infos = strings.Split(info, "\r\n")
	return
}
