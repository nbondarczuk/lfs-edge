// Represents configuration setings for the Krypton lfs-edge sync service.
package config

import "go.uber.org/zap"

// Server contains configuration settings for the Sync server.
// This includes the gRPC server and web sync server.
type Server struct {
	// Hostname of the sync service.
	Host string `yaml:"host"`
	// Port on which the gRPC server is available.
	RpcPort int `yaml:"rpc_port"`
}

// Database configuration settings
// embedded file based db
type Database struct {
	// database file path
	Path string `yaml:"path"`
	// schema
	Schema string `yaml:"schema"`
	// migrate
	Migrate bool `yaml:"migrate"`
}

// Files storage configuration settings
type Files struct {
	// Storage cache location path
	StoragePath string `yaml:"storage_path"`
	// Maximal disk usage in Gb
	MaxDiskUse string `yaml:"max_disk_use"`
	// Space management policy ie. delete_oldest, delete_never_used, etc.
	SpaceManagementPolicy string `yaml:"space_management_policy"`
	// Retry paramter returned to proxy
	RetryAfterSeconds int `yaml:"retry_after_seconds"`
	// FS server for remote file fetch  to be used in the template
	FileServerHost string `yaml:"file_server_host"`
	// FS server port to be used in the template
	FileServerPort int `yaml:"file_server_port"`
	// Template of API request to fetch the remorte files from fs
	FileServerURLTemplate string `yaml:"file_server_url_template"`
	// Size of buffered channel
	PendingChannelSize int  `yaml:"pending_channel_size"`
}

// Config represents configuration settings for the sync service.
type Config struct {
	// Configuration settings for the gRPC and REST servers.
	Server `yaml:"server"`
	// Database settings
	Database `yaml:"database"`
	// Files storage settings
	Files `yaml:"files"`
	// Whether the Sync services are configured to run in test mode.
	TestMode bool `yaml:"test_mode"`
	// command line switches and version info
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
	Logger *zap.Logger
}
