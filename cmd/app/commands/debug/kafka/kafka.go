package kafka

import (
	"app/internal/database"
	"app/internal/database/kafka"
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
	kafka.WithCobra(Command)
	Command.Flags().VisitAll(func(f *pflag.Flag) {
		viper.GetViper().BindPFlag(f.Name, f)
	})
}

var Command = &cobra.Command{
	Use:   "kafka",
	Short: "Debugs the Kafka connection",
	RunE:  run,
}

func run(command *cobra.Command, args []string) error {
	scheme := viper.GetString(kafka.ConfigScheme)
	logrus.Infof("scheme: '%s'", scheme)

	user := viper.GetString(kafka.ConfigUser)
	logrus.Infof("user: '%s'", user)

	password := viper.GetString(kafka.ConfigPassword)
	logrus.Infof("password: '%s'", password)

	hosts := viper.GetStringSlice(kafka.ConfigHost)
	logrus.Infof("hosts: [%s]", strings.Join(hosts, ","))

	options := viper.GetStringSlice(kafka.ConfigOption)
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

	tlsCaCertPath, err := command.Flags().GetString(kafka.ConfigCaCertPath)
	if err != nil {
		return fmt.Errorf("failed to get ca certificate path: %s", err)
	}
	logrus.Infof("tls ca cert path: %s", tlsCaCertPath)
	tlsClientCertPath, err := command.Flags().GetString(kafka.ConfigClientCertPath)
	if err != nil {
		return fmt.Errorf("failed to get client certificate path: %s", err)
	}
	logrus.Infof("tls client cert path: %s", tlsClientCertPath)
	tlsClientKeyPath, err := command.Flags().GetString(kafka.ConfigClientKeyPath)
	if err != nil {
		return fmt.Errorf("failed to get client key path: %s", err)
	}
	logrus.Infof("tls client key path: %s", tlsClientKeyPath)
	var tlsConfig *kafka.TlsOpts = nil
	if tlsCaCertPath != "" && tlsClientCertPath != "" && tlsClientKeyPath != "" {
		tlsConfig = &kafka.TlsOpts{
			CaCertificatePath:     tlsCaCertPath,
			ClientCertificatePath: tlsClientCertPath,
			ClientKeyPath:         tlsClientKeyPath,
		}
	}
	logrus.Infof("tls enabled: %v", tlsConfig != nil)

	connection, err := kafka.NewConnection(
		database.ConnectionOpts{
			Scheme:   scheme,
			Username: user,
			Password: password,
			Hosts:    hosts,
			Options:  optionsValue,
		},
		kafka.ConnectionOpts{
			TLS: tlsConfig,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to connect to kafka: %s", err)
	}
	broker := connection.Broker()
	logrus.Infof("successfully connected to kafka[%s] with id[%v] on port[%v] on rack[%s]", broker.Host, broker.ID, broker.Port, broker.Rack)
	connection.Close()
	return nil
}
