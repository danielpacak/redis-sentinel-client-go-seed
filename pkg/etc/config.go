package etc

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	RedisAddr   string `env:"SEED_REDIS_ADDR" envDefault:"redis.redis:26379"`
	RedisMaster string `env:"SEED_REDIS_MASTER" envDefault:"mymaster"`
}

func GetConfig() (cfg Config, err error) {
	err = env.Parse(&cfg)
	return
}
