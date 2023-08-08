package mongo

import (
	"app/internal/constants"
	"app/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

func NewConnection(opts database.ConnectionOpts) (*mongo.Client, error) {
	if err := initOpts(&opts); err != nil {
		return nil, err
	}
	connectionUri := fmt.Sprintf("%s://%s", opts.Scheme, strings.Join(opts.Hosts, ","))
	credentials := options.Credential{
		Username: opts.Username,
		Password: opts.Password,
	}
	if authSource, ok := opts.Options["authSource"]; ok && authSource != "" {
		credentials.AuthSource = authSource
	} else {
		credentials.AuthSource = opts.Database
	}
	connectOpts := options.Client().
		SetAppName(constants.AppName).
		ApplyURI(connectionUri).
		SetAuth(credentials).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetCompressors([]string{"zlib", "zstd"})

	if writeConcern, ok := opts.Options["writeConcern"]; ok && writeConcern != "" {
		writeConcernValue, err := strconv.Atoi(writeConcern)
		if err != nil {
			return nil, fmt.Errorf("failed to receive a valid writeConcern (received '%v') to connect to mongo via uri[%s]: %s", writeConcern, connectionUri, err)
		}
		journal := true
		connectOpts = connectOpts.SetWriteConcern(
			&writeconcern.WriteConcern{W: writeConcernValue, Journal: &journal},
		)
	}
	client, err := mongo.Connect(context.Background(), connectOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo via uri[%s]: %s", connectionUri, err)
	}
	return client, nil
}
