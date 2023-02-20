package config

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	MODULE                   = "lfs-edge proxy"
	COMPONENT                = "HP CEM lfs-edge proxy"	
	DEFAULT_CONFIG_FILE_NAME = "config.yml"
	DEFAULT_LOG_LEVEL        = "info"
)

var Settings Config

func Init() {
	initFlags()
	initLogger(DEFAULT_LOG_LEVEL)
}

func initFlags() {
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
	fmt.Printf("%s: version information", COMPONENT)
	fmt.Printf("- Git commit hash: %s\n - Built at: %s\n - Built by: %s\n - Built on: %s\n",
		Settings.Flags.gitCommitHash,
		Settings.Flags.builtAt,
		Settings.Flags.builtBy,
		Settings.Flags.builtOn)
}

func Load(testModeEnabled bool) bool {
	var filename string = DEFAULT_CONFIG_FILE_NAME

	// check required environment variables
	// assumes aws environment
	if !CheckRequiredEnvVars() {
		logger.Error("Could not find required environment variables!")
		return false
	}

	// Check if the default configuration file has been overridden using the
	// command line switch --config_file
	if *Settings.Flags.ConfigFile != "" {
		logger.Info("Using configuration file specified by command line switch.",
			zap.String("ConfigFile", *Settings.Flags.ConfigFile),
		)
		filename = *Settings.Flags.ConfigFile
	}

	// Open the configuration file for parsing.
	fh, err := os.Open(filename)
	if err != nil {
		logger.Error("Failed to load configuration file!",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return false
	}
	defer fh.Close()

	// Read the configuration file and unmarshal the YAML.
	decoder := yaml.NewDecoder(fh)
	err = decoder.Decode(&Settings)
	if err != nil {
		logger.Error("Failed to parse configuration file!",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return false
	}

	logger.Info("Parsed configuration from the configuration file!",
		zap.String("filename", filename),
	)

	// Override config from environment variables
	// note this only happens if environment variables are specified
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
	return logger
}
