package mysql

import (
	"github.com/spf13/cobra"
)

const (
	ConfigScheme   = "mysql-scheme"
	ConfigUser     = "mysql-user"
	ConfigPassword = "mysql-password"
	ConfigHost     = "mysql-host"
	ConfigDatabase = "mysql-database"
	ConfigOption   = "mysql-option"
)

func WithCobra(command *cobra.Command) {
	command.Flags().String(ConfigScheme, DefaultScheme, "Scheme used to connect to the database")
	command.Flags().String(ConfigUser, "user", "User used to authenticate with the database")
	command.Flags().String(ConfigPassword, "password", "Password used to authenticate with the database")
	command.Flags().StringSlice(ConfigHost, []string{"127.0.0.1:3306"}, "Hosts used to connect to the database")
	command.Flags().String(ConfigDatabase, "database", "Name of the database to connect to")
	command.Flags().StringSlice(ConfigOption, []string{"parseTime=true", "readTimeout=10s"}, "Connection options where applicable")
}
