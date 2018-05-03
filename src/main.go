package main

import (
	"fmt"
	"os/exec"
	"runtime"

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
	OpusBitrate     int
}

var (
	conf *Configuration
)

func readConf() *Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read Error: %v\n", err)
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
	systemCheck()
}

func systemCheck() {
	fmt.Println("--- WAV2LOSS PRE-RUN CHECK ---")
	var osVersion = runtime.GOOS

	fmt.Printf("I am running on: %s\n", osVersion)
	fmt.Printf("%+v\n", conf)

	opusBin, err := exec.LookPath("opusenc")
	if err != nil {
		fmt.Printf("Could not find opusenc in PATH\n")
	}
	lameBin, err := exec.LookPath("lame")
	if err != nil {
		fmt.Printf("Could not find LAME in PATH\n")
	}

	_ = opusBin
	_ = lameBin
	fmt.Printf("opus:\t%s\nLAME:\t%s\n", opusBin, lameBin)

}

// Main program
func main() {
	// TODO: Check config file parameters
	// TODO: Test executables

}
