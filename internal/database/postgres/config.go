package postgres

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "postgres-scheme"
	ConfigUser     = "postgres-user"
	ConfigPassword = "postgres-password"
	ConfigHost     = "postgres-host"
	ConfigDatabase = "postgres-database"
	ConfigOption   = "postgres-option"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, "user", "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, "password", "Password used to authenticate with the database")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:5432"}, "Hosts used to connect to the database")
	command.Flags().String(ConfigDatabase, "database", "Name of the database to connect to")
	command.Flags().StringSlice(ConfigOption, []string{}, "Connection options where applicable")
}
