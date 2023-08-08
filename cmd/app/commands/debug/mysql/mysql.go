package mysql

import (
	"app/internal/database"
	"app/internal/database/mysql"
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
	mysql.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "mysql",
	Short: "Debugs the MySQL connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(mysql.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(mysql.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(mysql.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(mysql.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	name := viper.GetString(mysql.ConfigDatabase)
	logrus.Infof("name: '%s'", name)

	options := viper.GetStringSlice(mysql.ConfigOption)
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

	connection, err := mysql.NewConnection(database.ConnectionOpts{
		Scheme:   scheme,
		Username: user,
		Password: password,
		Hosts:    hosts,
		Database: name,
		Options:  optionsValue,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %s", err)
	}
	if err := connection.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql: %s", err)
	}
	return nil
}
