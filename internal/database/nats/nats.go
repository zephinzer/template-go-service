package nats

import (
	"app/internal/database"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/sirupsen/logrus"
)

const (
	DefaultScheme = "nats"
)

func NewConnection(opts database.ConnectionOpts, natsOpts ...ConnectionOpts) (*nats.Conn, error) {
	natsOptions := []nats.Option{
		nats.PingInterval(500 * time.Millisecond),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			logrus.Warnf("nats error: %s", err)
		}),
	}
	isNkeyUsed := false
	nkeyValue := ""

	if len(natsOpts) > 0 {
		if natsOpts[0].NKey != "" {
			isNkeyUsed = true
			nkeyValue = natsOpts[0].NKey
		}
	}

	if isNkeyUsed {
		keyPair, err := nkeys.ParseDecoratedNKey([]byte(nkeyValue))
		if err != nil {
			return nil, fmt.Errorf("failed to parse provided nkey: %s", err)
		}
		defer keyPair.Wipe()
		publicKey, err := keyPair.PublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to receive a valid public key: %s", err)
		}
		if !nkeys.IsValidPublicUserKey(publicKey) {
			return nil, fmt.Errorf("failed to receive a valid public user key: %s", err)
		}
		natsOptions = append(natsOptions, nats.Nkey(publicKey, func(nonce []byte) ([]byte, error) {
			signature, err := keyPair.Sign([]byte(nonce))
			if err != nil {
				return nil, fmt.Errorf("failed to sign key with nonce[%s]: %s", string(nonce), err)
			}
			return signature, nil
		}))
		opts.Username = ""
		opts.Password = ""
	}
	natsUrl := newNatsUrl(opts.Scheme, opts.Username, opts.Password, opts.Hosts)
	connection, err := nats.Connect(natsUrl, natsOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats using scheme[%s] and hosts['%s']: %s", opts.Scheme, strings.Join(opts.Hosts, "','"), err)
	}
	return connection, nil
}

func newNatsUrl(
	scheme string,
	username string,
	password string,
	hosts []string,
) string {
	natsUrls := []string{}
	for _, host := range hosts {
		auth := ""
		if username != "" && password != "" {
			auth = username + ":" + password + "@"
		}
		natsUrls = append(natsUrls, scheme+"://"+auth+host)
	}
	return strings.Join(natsUrls, ",")
}
