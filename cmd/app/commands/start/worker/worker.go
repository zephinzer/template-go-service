package worker

import (
	"app/internal/example"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "worker",
	Short: "Starts the worker component",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	example.StartWorker()
	return nil
}
