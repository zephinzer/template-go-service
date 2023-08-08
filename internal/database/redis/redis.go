package redis

import (
	"app/internal/constants"
	"app/internal/database"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultScheme = "tcp"
)

func initOpts(opts *database.ConnectionOpts) error {
	issues := []error{}
	if opts.Password == "" {
		issues = append(issues, errors.New("failed to receive a password to use"))
	}
	if len(opts.Hosts) == 0 {
		issues = append(issues, errors.New("failed to receive at least one host to connect to"))
	}
	if opts.Scheme == "" {
		opts.Scheme = DefaultScheme
	}
	if len(issues) > 0 {
		return errors.Join(issues...)
	}
	return nil
}

func NewConnection(opts database.ConnectionOpts, redisOpts ...ConnectionOpts) (*redis.Client, error) {
	db := 0
	if opts.Database != "" {
		var err error
		db, err = strconv.Atoi(opts.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to select database[%v] (needs to be an integer): %s", opts.Database, err)
		}
	}
	redisOptions := redis.Options{
		ClientName:            constants.AppName,
		Network:               opts.Scheme,
		Addr:                  strings.Join(opts.Hosts, ","),
		Username:              opts.Username,
		Password:              opts.Password,
		DB:                    db,
		ContextTimeoutEnabled: true,

		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  8 * time.Second,

		ConnMaxIdleTime: 3 * time.Second,
		ConnMaxLifetime: 120 * time.Second,
		MinIdleConns:    10,
		MaxIdleConns:    50,
	}
	if len(redisOpts) > 0 {
		redisOptions.ConnMaxIdleTime = redisOpts[0].ConnMaxIdleTime
		redisOptions.ConnMaxLifetime = redisOpts[0].ConnMaxLifetime
		redisOptions.PoolFIFO = redisOpts[0].IsPoolFIFO
		redisOptions.MaxIdleConns = redisOpts[0].MaxIdleConns

		redisOptions.DialTimeout = redisOpts[0].DialTimeout
		redisOptions.ReadTimeout = redisOpts[0].ReadTimeout
		redisOptions.WriteTimeout = redisOpts[0].WriteTimeout
		redisOptions.PoolTimeout = redisOpts[0].PoolTimeout
	}
	cache := redis.NewClient(&redisOptions)

	return cache, nil
}
