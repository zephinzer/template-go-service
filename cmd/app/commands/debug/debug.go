package debug

import (
	"app/cmd/app/commands/debug/kafka"
	"app/cmd/app/commands/debug/mongo"
	"app/cmd/app/commands/debug/mysql"
	"app/cmd/app/commands/debug/nats"
	"app/cmd/app/commands/debug/postgres"
	"app/cmd/app/commands/debug/redis"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(kafka.Command)
	Command.AddCommand(mongo.Command)
	Command.AddCommand(mysql.Command)
	Command.AddCommand(nats.Command)
	Command.AddCommand(postgres.Command)
	Command.AddCommand(redis.Command)
}

var Command = &cobra.Command{
	Use:   "debug",
	Short: "Debug stuff",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	return command.Help()
}
