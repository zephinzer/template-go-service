package job

import (
	"app/internal/example"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().Duration("duration", 3*time.Second, "Specifies a timeout after which the job will complete")
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "job",
	Short: "Starts a job",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	duration := viper.GetDuration("duration")
	example.StartJob(duration)
	return nil
}
