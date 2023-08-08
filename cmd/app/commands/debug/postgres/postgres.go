package postgres

import (
	"app/internal/database"
	"app/internal/database/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	postgres.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "postgres",
	Short: "Debugs the PostgreSQL connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(postgres.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(postgres.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(postgres.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(postgres.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	name := viper.GetString(postgres.ConfigDatabase)
	logrus.Infof("name: '%s'", name)

	options := viper.GetStringSlice(postgres.ConfigOption)
	logrus.Infof("options: [%s]", strings.Join(options, ","))
	optionsValue := map[string]string{}
	for _, option := range options {
		optionKeyValue := strings.SplitN(option, "=", 2)
		if len(optionKeyValue) != 2 {
			return errors.New("option key[%s] does not have a value")
		}
		optionsValue[optionKeyValue[0]] = optionKeyValue[1]
	}
	optionsString, _ := json.Marshal(optionsValue)
	logrus.Infof("options map: %s", string(optionsString))

	connection, err := postgres.NewConnection(database.ConnectionOpts{
		Scheme:   scheme,
		Username: user,
		Password: password,
		Hosts:    hosts,
		Database: name,
		Options:  optionsValue,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %s", err)
	}
	pingTimeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := connection.Ping(pingTimeoutCtx); err != nil {
		return fmt.Errorf("failed to ping postgres: %s", err)
	}
	return nil
}
