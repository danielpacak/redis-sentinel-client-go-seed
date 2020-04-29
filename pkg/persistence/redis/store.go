package redis

import (
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

func (s *store) Keys() (keys []string, err error) {
	conn := s.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("KEYS", "*")
	if err != nil {
		return
	}
	log.Debugf("Reply: %v", reply)
	keys, err = xredis.Strings(reply, err)
	return
}
