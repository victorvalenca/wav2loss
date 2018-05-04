// +build amd64

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
	conf     *Configuration
	opusBin  string
	lameBin  string
	simulate = flag.Bool("simulate", false, "if set, no underlying commands will be run")
)

func readConf() *Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read in YAML file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read Error: %v\n", err)
	}

	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("Unable to parse configuration file into struct. %v", err)
	}

	// Parse potential envars
	switch conf.RecordDirectory[0] {
	case '$', '%':
		conf.RecordDirectory = os.Getenv(strings.Trim(conf.RecordDirectory, "$%"))
	}
	switch conf.OutputDirectory[0] {
	case '$', '%':
		conf.OutputDirectory = os.Getenv(strings.Trim(conf.OutputDirectory, "$%"))
	}

	return conf
}

// Initialization routine
func init() {
	conf = readConf()
	pass := systemCheck()
	if !pass {
		fmt.Printf("Configuration Check failed. Exiting...\n---\n")
		os.Exit(1)
	}
}

func systemCheck() bool {
	fmt.Println("--- WAV2LOSS PRE-RUN CHECK ---")
	var osVersion = runtime.GOOS

	fmt.Printf("I am running on: %s\n", osVersion)

	// Look for executables for opus and lame
	res, err := exec.LookPath("opusenc")
	if err != nil {
		fmt.Printf("Could not find opusenc in PATH\n")
		return false
	}
	opusBin = res

	res, err = exec.LookPath("lame")
	if err != nil {
		fmt.Printf("Could not find LAME in PATH\n")
		return false
	}
	lameBin = res
	fmt.Printf("Configuration Check passed\n---\n")
	return true
}

// Main program
func main() {

	if len(os.Args) < 2 {
		fmt.Printf("No File Given. Exiting...\n" +
			"USAGE: wav2loss filename_in_recording_directory.wav")
		os.Exit(1)
	}

	// Check for "simulation" mode flag
	flag.Parse()

	// opus/lame args: [options] input output
	t := time.Now()
	tFormatted := t.UTC().Format("2006-01-02")
	trimFile := strings.Replace(conf.Title, " ", "_", -1)
	outFile := filepath.Join(conf.OutputDirectory, trimFile+"_"+tFormatted)
	inFile := filepath.Join(conf.RecordDirectory, os.Args[1])

	opusTest := exec.Command(opusBin,
		"--bitrate", conf.OpusBitrate,
		"--title", "\""+conf.Title+"\"",
		"--artist", "\""+conf.Artist+"\"",
		"--album", "\""+conf.Album+"\"",
		"--date", "\""+t.UTC().Format("2006")+"\"",
		inFile,
		outFile+".opus")

	lameTest := exec.Command(lameBin,
		"-"+conf.LameBitrate,
		"--add-id3v2",
		"--tt", "\""+conf.Title+"\"",
		"--ta", "\""+conf.Artist+"\"",
		"--tl", "\""+conf.Album+"\"",
		"--ty", "\""+tFormatted+"\"",
		inFile,
		outFile+".mp3")

	opusTest.Stdout, lameTest.Stdout = os.Stdout, os.Stdout
	opusTest.Stderr, lameTest.Stderr = os.Stderr, os.Stderr

	if !*simulate {
		err := opusTest.Run()
		if err != nil {
			log.Printf("wut, couldn't run opusenc :( %v\n", err)
		}
		err = lameTest.Run()
		if err != nil {
			log.Printf("wut, couldn't run lame :( %v\n", err)
		}

	} else {
		fmt.Printf("SIMULATION MODE, PRINTING ARGUMENTS INSTEAD\n")
		fmt.Printf("OPUSENC ARGS: %v\n", opusTest.Args)
		fmt.Printf("LAME ARGS: %v\n", lameTest.Args)
	}
}
