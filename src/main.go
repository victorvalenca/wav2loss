package main

import (
	"fmt"
	"os"
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
	conf    *Configuration
	opusBin string
	lameBin string
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

	res, err := exec.LookPath("opusenc")
	if err != nil {
		fmt.Printf("Could not find opusenc in PATH\n")
	}
	opusBin = res

	res2, err := exec.LookPath("lame")
	if err != nil {
		fmt.Printf("Could not find LAME in PATH\n")
	}
	lameBin = res2
}

// Configuration check routine
func checkConf() bool {

	// TODO: Parse envvars to full path
	switch conf.RecordDirectory[0] {
	case '$':
		fmt.Printf("I am dealing with a UNIX-style envvar, use full path please\n")
		return false
	case '%':
		fmt.Printf("I am dealing with a Windows-style envvar, use full path please\n")
		return false
	default:
		return true
	}
}

// Main program
func main() {
	// TODO: Test executables

	fmt.Printf("opus:\t%s\nLAME:\t%s\n", opusBin, lameBin)
	opusTest := exec.Command(opusBin, "--version")
	lameTest := exec.Command(lameBin, "--version")

	opusTest.Stdout, lameTest.Stdout = os.Stdout, os.Stdout
	opusTest.Stderr, lameTest.Stderr = os.Stderr, os.Stderr

	pass := checkConf()
	if !pass {
		fmt.Printf("Configuration check failed\n")
	}
	// opusTest.Run()
	// lameTest.Run()

}
