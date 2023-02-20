package config

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type value struct {
	secret bool
	v      interface{}
}

// check for environment variables that are required for operation
func CheckRequiredEnvVars() bool {
	var isPresent bool
	keys := [...]string{}
	for _, key := range keys {
		_, isPresent = os.LookupEnv(key)
		if !isPresent {
			syncLogger.Info("Could not find required environment variable.",
				zap.String("variable: ", key))
			return false
		}
	}
	return true
}

// OverrideFromEnvironment - check values specified for supported
// environment variables. These can be used to override configuration settings
// specified in the config file.
func (c *Config) OverrideFromEnvironment() {
	m := map[string]value{
		// Server
		"LFSSYNC_SERVER_HOST":                    {v: &c.Server.Host},
		"LFSSYNC_SERVER_RPC_PORT":                {v: &c.Server.RpcPort},
		// Database
		"LFSSYNC_DATABASE_PATH":                  {v: &c.Database.Path},
		"LFSSYNC_DATABASE_SCHEMA":                {v: &c.Database.Schema},		
		"LFSSYNC_DATABASE_MIGRATE":               {v: &c.Database.Migrate},
		// Files
		"LFSSYNC_FILES_STORAGE_PATH":             {v: &c.Files.StoragePath},
		"LFSSYNC_FILES_MAX_DISK_USE":             {v: &c.Files.MaxDiskUse},
		"LFSSYNC_FILES_SPACE_MANAGEMENT_POLICY":  {v: &c.Files.SpaceManagementPolicy},
		"LFSSYNC_FILES_FILE_SERVER_HOST":         {v: &c.Files.FileServerHost},
		"LFSSYNC_FILES_FILE_SERVER_PORT":         {v: &c.Files.FileServerPort},
		"LFSSYNC_FILES_FILE_SERVER_URL_TEMPLATE": {v: &c.Files.FileServerURLTemplate},
		"LFSSYNC_FILES_PENDING_CHANNEL_SZIE":     {v: &c.Files.PendingChannelSize},				
	}
	for k, v := range m {
		e := os.Getenv(k)
		if e != "" {
			syncLogger.Info("Overriding configuration from environment variable.",
				zap.String("variable: ", k),
				zap.String("value: ", getLoggableValue(v.secret, e)))
			replaceConfigValue(os.Getenv(k), &v)
		}
	}
}

// envValue will be non empty as this function is private to file
func replaceConfigValue(envValue string, t *value) {
	switch t.v.(type) {
	case *string:
		*t.v.(*string) = envValue
	case *[]string:
		valSlice := strings.Split(envValue, ",")
		for i := range valSlice {
			valSlice[i] = strings.TrimSpace(valSlice[i])
		}
		*t.v.(*[]string) = valSlice
	case *bool:
		b, err := strconv.ParseBool(envValue)
		if err != nil {
			syncLogger.Error("Bad bool value in env")
		} else {
			*t.v.(*bool) = b
		}
	case *int:
		i, err := strconv.Atoi(envValue)
		if err != nil {
			syncLogger.Error("Bad integer value in env",
				zap.Error(err))
		} else {
			*t.v.(*int) = i
		}
	default:
		syncLogger.Error("There was a bad type map in env override",
			zap.String("value", envValue))
	}
}

func getLoggableValue(secret bool, value string) string {
	if secret {
		return "***"
	}
	return value
}
