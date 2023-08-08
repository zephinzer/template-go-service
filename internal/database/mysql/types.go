package mysql

import "time"

type ConnectionOpts struct {
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
