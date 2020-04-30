package etc

import (
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	RedisAddr   string `env:"SEED_REDIS_ADDR" envDefault:"redis.redis:26379"`
	RedisMaster string `env:"SEED_REDIS_MASTER" envDefault:"mymaster"`

	RedisConnectTimeout time.Duration `env:"SEED_REDIS_CONNECT_TIMEOUT" envDefault:"500ms"`
	RedisReadTimeout    time.Duration `env:"SEED_REDIS_READ_TIMEOUT" envDefault:"500ms"`
	RedisWriteTimeout   time.Duration `env:"SEED_REDIS_WRITE_TIMEOUT" envDefault:"500ms"`

	RedisPoolMaxIdle     int           `env:"SEED_REDIS_POOL_MAX_IDLE" envDefault:"3"`
	RedisPoolMaxActive   int           `env:"SEED_REDIS_POOL_MAX_ACTIVE" envDefault:"10"`
	RedisPoolIdleTimeout time.Duration `env:"SEED_REDIS_POOL_IDLE_TIMEOUT" envDefault:"5m"`
}

func GetConfig() (cfg Config, err error) {
	err = env.Parse(&cfg)
	return
}

func GetLogLevel() log.Level {
	if value, ok := os.LookupEnv("SEED_LOG_LEVEL"); ok {
		level, err := log.ParseLevel(value)
		if err != nil {
			return log.InfoLevel
		}
		return level
	}
	return log.InfoLevel
}
