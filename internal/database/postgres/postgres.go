package postgres

import (
	"app/internal/database"
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

const (
	DefaultScheme = "postgres"
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

func NewConnection(opts database.ConnectionOpts) (*pgx.Conn, error) {
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
		"%s://%s:%s@%s/%s%s",
		opts.Scheme,
		opts.Username,
		opts.Password,
		strings.Join(opts.Hosts, ","),
		opts.Database,
		encodedOptions,
	)
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to postgres: %s", err)
	}

	return db, nil
}
