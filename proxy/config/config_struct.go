package config

import (
	"go.uber.org/zap"
)

type TokenType string

// Configuration settings for the REST server.
type Server struct {
	Host string `yaml:"host"`

	// Port on which the REST service is available.
	Port int `yaml:"port"`

	// Debug rest requests
	DebugRestRequests bool `yaml:"debug_rest_requests"`
}

type RpcServer struct {
	// Hostname of the sync gRPC server
	Host string `yaml:"host"`
	// Port of the sync gRPC server
	Port int `yaml:"port"`
}

// Configuration settings for files.
type Files struct {
	Path string `yaml:"path"`
}

// Config structire for all settings
type Config struct {
	// Rest server settings
	Server Server  `yaml:"server"`

	// gRPC server setup for client
	RpcServer RpcServer `yaml:"rpc_server"`
	
	// Files settings
	Files Files  `yaml:"files"`

	// Command line switches/flags.d
	Flags struct {
		// --config_file: specifies the path to the configuration file.
		ConfigFile *string
		// --log_level: specify the logging level to use.
		LogLevel *string
		// --version: displays versioning information.
		Version *bool
		//
		gitCommitHash string
		builtAt       string
		builtBy       string
		builtOn       string
	}

	// Structured logging using Uber Zap.
	Logger *zap.Logger

	// Whether the service is running in test mode.
	TestMode bool
}
