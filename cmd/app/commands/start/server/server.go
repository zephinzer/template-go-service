package server

import (
	"app/internal/api"
	"app/internal/server"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	server.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP server",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	if err := server.StartHttp(server.StartHttpOpts{
		BindInterface: viper.GetString(server.ConfigInterface),
		BindPort:      viper.GetInt(server.ConfigPort),
		Router:        api.WithFiber,
	}); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}
	return nil
}
