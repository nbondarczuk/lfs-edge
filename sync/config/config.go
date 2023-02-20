package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	MODULE                   = "lfs-edge sync"
	COMPONENT                = "HP CEM lfs-edge sync"
	DEFAULT_CONFIG_FILE_NAME = "config.yml"
	DEFAULT_LOG_LEVEL        = "info"
)

var Settings Config

func Init() {
	initFlags()
	initLogger(DEFAULT_LOG_LEVEL)
}

func initFlags() {
	// Set up required command line flags.
	Settings.Flags.ConfigFile = flag.String("config_file", "",
		"Specify the path to the configuration file.")
	Settings.Flags.LogLevel = flag.String("log_level", "", "Specify the logging level.")
	Settings.Flags.Version = flag.Bool("version", false,
		"Print the version of the service and exit!")

	// Parse the command line flags.
	flag.Parse()
	if *Settings.Flags.Version {
		printVersionInformation()
	}
}

func printVersionInformation() {
	fmt.Printf("%s: version information\n", COMPONENT)
	fmt.Printf("- Git commit hash: %s\n - Built at: %s\n - Built by: %s\n - Built on: %s\n",
		Settings.Flags.gitCommitHash,
		Settings.Flags.builtAt,
		Settings.Flags.builtBy,
		Settings.Flags.builtOn)
}

// Load acquires config from YAML file and overrides with env variables.
// TEST_MODE env variable triggers test mode.
func Load(testModeEnabled bool) bool {
	var filename string = DEFAULT_CONFIG_FILE_NAME

	// Check required environment variables.
	// It assumes AWS environment.
	if !CheckRequiredEnvVars() {
		syncLogger.Error("Could not find required environment variables!")
		return false
	}

	fmt.Println(*Settings.Flags.ConfigFile)
	// Check if the default configuration file has been overridden using the
	// command line switch --config_file
	if *Settings.Flags.ConfigFile != "" {
		syncLogger.Info("Using configuration file specified by command line switch.",
			zap.String("Configuration file:", *Settings.Flags.ConfigFile),
		)
		filename = *Settings.Flags.ConfigFile
	}

	// Open the configuration file for parsing.
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		syncLogger.Error("Failed to load configuration file!",
			zap.String("Configuration file:", filename),
			zap.Error(err),
		)
		return false
	}

	// Read the configuration file and unmarshal the YAML.
	err = yaml.Unmarshal(bytes, &Settings)
	if err != nil {
		syncLogger.Error("Failed to parse configuration file!",
			zap.String("Configuration file:", filename),
			zap.Error(err),
		)
		return false
	}

	syncLogger.Info("Parsed configuration from the configuration file!",
		zap.String("Configuration file:", filename),
	)

	// Override config from environment variables.
	// Note: This only happens if environment variables are specified.
	Settings.OverrideFromEnvironment()

	testModeEnvVar := os.Getenv("TEST_MODE")
	if (testModeEnvVar == "enabled") || (testModeEnabled) {
		Settings.TestMode = true
		fmt.Printf("%s is running in test mode with test hooks enabled.\n", MODULE)
		InitTestLogger()
	}

	return true
}

func Shutdown() {
	shutdownLogger()
}

func GetLogger() *zap.Logger {
	return syncLogger
}
