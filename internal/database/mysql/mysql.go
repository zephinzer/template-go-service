package mysql

import (
	"app/internal/database"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DefaultScheme = "tcp"
)

func initOpts(opts *database.ConnectionOpts) error {
	issues := []error{}
	if opts.Database == "" {
		issues = append(issues, errors.New("failed to receive a database to connect to"))
	}
	if opts.Username == "" {
		issues = append(issues, errors.New("failed to receive a username to use"))
	}
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

func NewConnection(opts database.ConnectionOpts, mysqlOpts ...ConnectionOpts) (*sql.DB, error) {
	if err := initOpts(&opts); err != nil {
		return nil, err
	}
	hasAtLeastOneOption := false
	options := url.Values{}
	for key, value := range opts.Options {
		options.Add(key, value)
		hasAtLeastOneOption = true
	}
	encodedOptions := ""
	if hasAtLeastOneOption {
		encodedOptions = "?" + options.Encode()
	}
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s)/%s%s",
		opts.Username,
		opts.Password,
		opts.Scheme,
		strings.Join(opts.Hosts, ","),
		opts.Database,
		encodedOptions,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to mysql: %s", err)
	}
	connMaxLifetime := time.Minute * 3
	maxIdleConns := 10
	maxOpenConns := 10
	if len(mysqlOpts) > 0 {
		if mysqlOpts[0].ConnMaxLifetime > 0 {
			connMaxLifetime = mysqlOpts[0].ConnMaxLifetime
		}
		if mysqlOpts[0].MaxIdleConns > 0 {
			maxIdleConns = mysqlOpts[0].MaxIdleConns
		}
		if mysqlOpts[0].MaxOpenConns > 0 {
			maxOpenConns = mysqlOpts[0].MaxOpenConns
		}
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxIdleConns)
	db.SetMaxIdleConns(maxOpenConns)

	return db, nil
}
