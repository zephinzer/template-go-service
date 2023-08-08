package redis

import "time"

type ConnectionOpts struct {
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	IsPoolFIFO      bool
	MaxIdleConns    int

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
}
