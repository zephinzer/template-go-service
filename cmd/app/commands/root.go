package commands

import (
	"app/cmd/app/commands/debug"
	"app/cmd/app/commands/start"
	"app/internal/constants"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	Root.AddCommand(debug.Command)
	Root.AddCommand(start.Command)
	Root.Version = constants.Version
}

var Root = &cobra.Command{
	Use:   constants.AppName,
	Short: "Example CLI structure",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	return command.Help()
}
