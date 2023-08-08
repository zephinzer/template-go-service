package mongo

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "mongo-scheme"
	ConfigUser     = "mongo-user"
	ConfigPassword = "mongo-password"
	ConfigHost     = "mongo-host"
	ConfigDatabase = "mongo-database"
	ConfigOption   = "mongo-option"

	DefaultScheme   = "mongodb"
	DefaultUser     = "user"
	DefaultPassword = "password"
	DefaultDatabase = "database"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, DefaultUser, "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, DefaultPassword, "Password used to authenticate with the database")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:27017"}, "Hosts used to connect to the database")
	command.Flags().String(ConfigDatabase, DefaultDatabase, "Name of the database to connect to")
	command.Flags().StringSlice(ConfigOption, []string{"authSource=admin"}, "Connection options where applicable")
}
