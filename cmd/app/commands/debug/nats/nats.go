package nats

import (
	"app/internal/database"
	"app/internal/database/nats"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	nats.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "nats",
	Short: "Debugs the NATS connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(nats.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(nats.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(nats.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(nats.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	options := viper.GetStringSlice(nats.ConfigOption)
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

	nkey, err := command.Flags().GetString(nats.ConfigNkey)
	if err != nil {
		return fmt.Errorf("failed to get nkey: %s", err)
	}
	connection, err := nats.NewConnection(
		database.ConnectionOpts{
			Scheme:   scheme,
			Username: user,
			Password: password,
			Hosts:    hosts,
			Options:  optionsValue,
		},
		nats.ConnectionOpts{
			NKey: nkey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create connection to nats: %s", err)
	}
	status := connection.Status()
	if status.String() != "CONNECTED" {
		return fmt.Errorf("failed to connect to nats: %s", err)
	}
	serverName := connection.ConnectedServerName()
	logrus.Infof("connected to server[%s]", serverName)
	return nil
}
