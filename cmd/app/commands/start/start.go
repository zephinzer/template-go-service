package start

import (
	"app/cmd/app/commands/start/job"
	"app/cmd/app/commands/start/server"
	"app/cmd/app/commands/start/worker"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(job.Command)
	Command.AddCommand(server.Command)
	Command.AddCommand(worker.Command)
}

var Command = &cobra.Command{
	Use:   "start",
	Short: "Starts stuff",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	return command.Help()
}
