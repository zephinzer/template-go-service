package nats

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "nats-scheme"
	ConfigUser     = "nats-user"
	ConfigPassword = "nats-password"
	ConfigNkey     = "nats-nkey"
	ConfigHost     = "nats-host"
	ConfigOption   = "nats-option"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, "user", "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, "password", "Password used to authenticate with the database")
	command.Flags().String(ConfigNkey, "", "NKey used to authenticate with NATS")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:4222"}, "Hosts used to connect to the database")
	command.Flags().StringSlice(ConfigOption, []string{}, "Connection options where applicable")
}
