package redis

import (
	"app/internal/database"
	"app/internal/database/redis"
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
	redis.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "redis",
	Short: "Debugs the Redis connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(redis.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(redis.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(redis.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(redis.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	name := command.Flag(redis.ConfigDatabase).Value.String()
	logrus.Infof("name: '%s'", name)

	options := viper.GetStringSlice(redis.ConfigOption)
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

	connection, err := redis.NewConnection(database.ConnectionOpts{
		Scheme:   scheme,
		Username: user,
		Password: password,
		Hosts:    hosts,
		Database: name,
		Options:  optionsValue,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %s", err)
	}
	pingTimeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if status := connection.Ping(pingTimeoutCtx); status.Err() != nil {
		return fmt.Errorf("failed to ping redis: %s", status.Err())
	}
	return nil
}
