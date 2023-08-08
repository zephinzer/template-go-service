package redis

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "redis-scheme"
	ConfigUser     = "redis-user"
	ConfigPassword = "redis-password"
	ConfigHost     = "redis-host"
	ConfigDatabase = "redis-database"
	ConfigOption   = "redis-option"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, "user", "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, "password", "Password used to authenticate with the database")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:6379"}, "Hosts used to connect to the database")
	command.Flags().String(ConfigDatabase, "0", "ID of the database to connect to")
	command.Flags().StringSlice(ConfigOption, []string{}, "Connection options where applicable")
}
