package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration struct
type Configuration struct {
	Album           string
	Artist          string
	Title           string
	RecordDirectory string
	OutputDirectory string
	LameBitrate     string
	OpusBitrate     string
}

var (
	conf *Configuration
)

func readConf() *Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config.conf")
	viper.SetConfigType("conf")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("Unable to parse configuration file into struct. %v", err)
	}

	return conf
}

// Initialization routine
func init() {
	conf = readConf()
}

// Main program
func main() {
	fmt.Printf("%+v\n", conf)
}
