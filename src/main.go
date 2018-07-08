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
	conf      *Configuration
	opusBin   string
	lameBin   string
	simulate  = flag.Bool("simulate", false, "if set, no underlying commands will be run, and prints the command string instead")
	directIn  = flag.String("in", "", "direct path to the input file")
	directOut = flag.String("out", "", "direct path to the output directory")
)

func readConf() *Configuration {
	// Get the path to the running executable so viper knows where to find the config file
	here, err := os.Executable()
	viper.AddConfigPath(filepath.Dir(here))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read in YAML file
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read Error: %v\n", err)
		os.Exit(1)
	}

	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("Unable to parse configuration file into struct. %v", err)
		os.Exit(1)
	}

	// Parse potential envars
	if *directIn == "" {
		switch conf.RecordDirectory[0] {
		case '$', '%':
			if pathExists(conf.RecordDirectory) {
				conf.RecordDirectory = os.Getenv(strings.Trim(conf.RecordDirectory, "$%"))
			}
		}
	} else {
		conf.RecordDirectory = *directIn
	}
	if *directOut == "" {
		switch conf.OutputDirectory[0] {
		case '$', '%':
			if pathExists(conf.RecordDirectory) {
				conf.OutputDirectory = os.Getenv(strings.Trim(conf.OutputDirectory, "$%"))
			}
		}
	} else {
		conf.OutputDirectory = *directOut
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

// A check function to determine if a path is valid
func pathExists(pathName string) bool {
	_, err := os.Stat(pathName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
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
	flag.Parse()

	var args = flag.Args()
	if len(args) < 1 && len(*directIn) < 1 {
		fmt.Printf("No File(s) Given. Exiting...\nUSAGE:\n\twav2loss [-simulate] filename_in_recording_directory.wav\nOr with direct filepaths:\n\twav2loss [-simulate] [-in='path/to/input/file'] [-out='path/to/output/directory']")
		os.Exit(1)
	}

	// opus/lame args: [options] input output
	t := time.Now()
	tFormatted := t.UTC().Format("2006-01-02")
	trimFile := strings.Replace(conf.Title, " ", "_", -1)
	inFile := *directIn
	outFile := filepath.Join(*directOut, trimFile+"_"+tFormatted)
	if *directIn == "" {
		inFile = filepath.Join(conf.RecordDirectory, args[0])
	}
	if *directOut == "" {
		outFile = filepath.Join(conf.OutputDirectory, trimFile+"_"+tFormatted)
	}

	opusTest := exec.Command(opusBin,
		"--bitrate", conf.OpusBitrate,
		"--title", conf.Title,
		"--artist", conf.Artist,
		"--album", conf.Album,
		"--date", t.UTC().Format("2006"),
		inFile,
		outFile+".opus")

	lameTest := exec.Command(lameBin,
		"-"+conf.LameBitrate,
		"--add-id3v2",
		"--tt", conf.Title,
		"--ta", conf.Artist,
		"--tl", conf.Album,
		"--ty", tFormatted,
		inFile,
		outFile+".mp3")

	opusTest.Stdout, lameTest.Stdout = os.Stdout, os.Stdout
	opusTest.Stderr, lameTest.Stderr = os.Stderr, os.Stderr

	// Check for "simulation" mode flag, otherwise execute in sequence
	if !*simulate {
		// Check if WAV filename exists
		if !pathExists(filepath.Dir(conf.RecordDirectory)) {
			fmt.Printf("The file '%s' does not exist in the recording directory.\n", args[0])
			os.Exit(1)
		}
		err := opusTest.Run()
		if err != nil {
			log.Printf("I couldn't run opusenc :(\n%v\n", err)
		}
		err = lameTest.Run()
		if err != nil {
			log.Printf("I couldn't run lame :(\n%v\n", err)
		}

	} else {
		fmt.Printf("SIMULATION MODE, PRINTING ARGUMENTS INSTEAD\n")
		fmt.Printf("OPUSENC ARGS: %v\n", opusTest.Args)
		fmt.Printf("LAME ARGS: %v\n", lameTest.Args)
	}
}
