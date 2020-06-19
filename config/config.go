package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Schema struct {
	AMQP struct {
		URL             string `mapstructure:"url"`
		Host            string `mapstructure:"host"`
		Port            string `mapstructure:"port"`
		Vhost           string `mapstructure:"vhost"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		ExchangeName    string `mapstructure:"exchange_name"`
		ExchangeType    string `mapstructure:"exchange_type"`
		QueueName       string `mapstructure:"queue_name"`
		ConsumerThreads int    `mapstructure:"consumer_threads"`
	} `mapstructure:"amqp"`

	MongoDB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Env      string `mapstructure:"env"`
		Replica  string `mapstructure:"replica"`
	} `mapstructure:"mongodb"`
}

var Config Schema

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")          // Look for config in current directory
	config.AddConfigPath("config/")    // Optionally look for config in the working directory.
	config.AddConfigPath("../config/") // Look for config needed for tests.
	config.AddConfigPath("../")        // Look for config needed for tests.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// fmt.Printf("Current Config: %+v", Config)
}
