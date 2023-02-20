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
			logger.Info("Could not find required environment variable.",
				zap.String("variable: ", key))
			return false
		}
	}
	return true
}

// loadEnvironmentVariableOverrides - check values specified for supported
// environment variables. These can be used to override configuration settings
// specified in the config file.
func (c *Config) OverrideFromEnvironment() {
	m := map[string]value{
		// REST Server
		"LFSPROXY_HOST": {v: &c.Server.Host},
		"LFSPROXY_PORT": {v: &c.Server.Port},

		// gRPC Server for client to connect to
		"LFSPROXY_RPC_HOST": {v: &c.RpcServer.Host},
		"LFSPROXY_RPC_PORT": {v: &c.RpcServer.Port},

		// Files configuration settings
		"LFSPROXY_FILES_PATH": {v: &c.Files.Path},
	}
	for k, v := range m {
		e := os.Getenv(k)
		if e != "" {
			logger.Info("Overriding configuration from environment variable.",
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
			logger.Error("Bad bool value in env")
		} else {
			*t.v.(*bool) = b
		}
	case *int:
		i, err := strconv.Atoi(envValue)
		if err != nil {
			logger.Error("Bad integer value in env",
				zap.Error(err))
		} else {
			*t.v.(*int) = i
		}
	default:
		logger.Error("There was a bad type map in env override",
			zap.String("value", envValue))
	}
}

func getLoggableValue(secret bool, value string) string {
	if secret {
		return "***"
	}
	return value
}
