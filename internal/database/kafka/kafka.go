package kafka

import (
	"app/internal/constants"
	"app/internal/database"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func NewConnection(opts database.ConnectionOpts, kafkaOpts ...ConnectionOpts) (*kafka.Conn, error) {
	// simple authentication and security layer (sasl) setup
	mechanism, err := scram.Mechanism(scram.SHA512, opts.Username, opts.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to setup authentication mechanism: %s", err)
	}

	// secure sockets layer (ssl) setup if applicable
	var tlsConfig *tls.Config = nil
	isTlsEnabled := false
	if len(kafkaOpts) > 0 {
		extraOpts := kafkaOpts[0]
		if extraOpts.TLS != nil {
			var err error
			if tlsConfig, err = newTLSConfig(*extraOpts.TLS); err != nil {
				return nil, fmt.Errorf("failed to setup tls configuration: %s", err)
			}
			isTlsEnabled = true
		}
	}

	dialer := &kafka.Dialer{
		ClientID:      constants.AppName,
		Timeout:       4 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: mechanism,
	}
	scheme := opts.Scheme
	hosts := strings.Join(opts.Hosts, ",")
	dialCtx, cancelDial := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDial()
	connection, err := dialer.DialContext(dialCtx, scheme, hosts)
	if err != nil {
		return nil, fmt.Errorf("failed to dial kafka via scheme[%s] to host[%s] (tls enabled: %v): %s", scheme, hosts, isTlsEnabled, err)
	}
	return connection, nil
}

func newTLSConfig(opts TlsOpts) (*tls.Config, error) {
	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
	}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(opts.ClientCertificatePath, opts.ClientKeyPath)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(opts.CaCertificatePath)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}
