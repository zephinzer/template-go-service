package server

import "github.com/spf13/cobra"

const (
	ConfigInterface = "server-interface"
	ConfigPort      = "server-port"
)

func WithCobra(command *cobra.Command) {
	command.Flags().StringP(ConfigInterface, "i", "0.0.0.0", "interface for server to listen on")
	command.Flags().IntP(ConfigPort, "p", 3000, "port for server to listen on")

}
