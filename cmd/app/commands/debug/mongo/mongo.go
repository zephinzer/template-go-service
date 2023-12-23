package mongo

import (
	"app/internal/database"
	"app/internal/database/mongo"
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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	mongo.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "mongo",
	Short: "Debugs the mongo connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(mongo.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(mongo.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(mongo.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(mongo.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	name := viper.GetString(mongo.ConfigDatabase)
	logrus.Infof("name: '%s'", name)

	options := viper.GetStringSlice(mongo.ConfigOption)
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

	connection, err := mongo.NewConnection(database.ConnectionOpts{
		Scheme:   scheme,
		Username: user,
		Password: password,
		Hosts:    hosts,
		Database: name,
		Options:  optionsValue,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to mongo: %s", err)
	}
	pingTimeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := connection.Ping(pingTimeoutCtx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping mongo: %s", err)
	}
	disconnectCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := connection.Disconnect(disconnectCtx); err != nil {
		return fmt.Errorf("failed to disconnect from mongo: %s", err)
	}
	return nil
}
