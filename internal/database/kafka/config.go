package kafka

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "kafka-scheme"
	ConfigUser     = "kafka-user"
	ConfigPassword = "kafka-password"
	ConfigHost     = "kafka-host"
	ConfigOption   = "kafka-option"

	ConfigCaCertPath     = "kafka-ca-cert-path"
	ConfigClientCertPath = "kafka-client-cert-path"
	ConfigClientKeyPath  = "kafka-client-key-path"

	DefaultScheme = "tcp"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, "user", "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, "password", "Password used to authenticate with the database")
	command.Flags().String(ConfigCaCertPath, "./.data/kafka/config/certs/ca-cert", "Path to Certificate Authority certificate")
	command.Flags().String(ConfigClientCertPath, "./.data/kafka/config/certs/client-cert", "Path to client certificate")
	command.Flags().String(ConfigClientKeyPath, "./.data/kafka/config/certs/client-key", "Path to client key")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:9092"}, "Hosts used to connect to the database")
	command.Flags().StringSlice(ConfigOption, []string{}, "Connection options where applicable")
}
