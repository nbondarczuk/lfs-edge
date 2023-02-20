package config

import (
	"os"
	"strconv"
	"testing"
)

func init() {
	InitTestLogger()
}

func TestDoesNotOverrideWhenNoEnvVariablePresent(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	expected := c.Server.Host
	c.OverrideFromEnvironment()
	if c.Server.Host != "" {
		t.Fatalf(
			`%s:, Expected Server.Host = %s,  found: %s`,
			issue, expected, c.Server.Host)
	}
}

func TestOverridesAllStringEnvVariables(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	type test struct {
		envVarName    string
		fieldVarName  string
		configVar     *string
		inputValue    string
		expectedValue string
	}
	tests := []test{
		{
			"LFSSYNC_SERVER_HOST", "Server.Host",
			&c.Server.Host, "localhost", "localhost",
		},
		{
			"LFSSYNC_DATABASE_PATH", "Database.Path",
			&c.Database.Path, "whatever-path", "whatever-path",
		},

		{
			"LFSSYNC_DATABASE_SCHEMA", "Database.Schema",
			&c.Database.Schema, "whatever-schema", "whatever-schema",
		},
		{
			"LFSSYNC_FILES_STORAGE_PATH", "Files.StoragePath",
			&c.Files.StoragePath, "whatever-storage", "whatever-storage",
		},
		{
			"LFSSYNC_FILES_MAX_DISK_USE", "Files.MaxDiskUse",
			&c.Files.MaxDiskUse, "whatever-use", "whatever-use",
		},
		{
			"LFSSYNC_FILES_SPACE_MANAGEMENT_POLICY", "Files.SpaceManagementPolicy",
			&c.Files.SpaceManagementPolicy, "whatever-policy", "whatever-policy",
		},
		{
			"LFSSYNC_FILES_FILE_SERVER_HOST", "Files.FileServerHost",
			&c.Files.FileServerHost, "whatever-host", "whatever-host",
		},
		{
			"LFSSYNC_FILES_FILE_SERVER_URL_TEMPLATE", "Files.FileServerURLTemplate",
			&c.Files.FileServerURLTemplate, "whatever-template", "whatever-template",
		},
	}

	for _, tc := range tests {
		os.Setenv(tc.envVarName, tc.inputValue)
		c.OverrideFromEnvironment()
		if *tc.configVar != tc.expectedValue {
			t.Fatalf(
				`%s:, Expected %s = %s, found: %s`,
				issue, tc.fieldVarName, tc.expectedValue, *tc.configVar,
			)
		}
	}
}

func TestOverridesIntEnvVariable(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	expected := 123
	os.Setenv("LFSSYNC_SERVER_RPC_PORT", strconv.Itoa(expected))
	c.OverrideFromEnvironment()
	if c.Server.RpcPort != expected {
		t.Fatalf(
			`%s:, Expected Server.RpcPort = %d, found: %d`,
			issue, expected, c.Server.RpcPort)
	}
}

func TestOverridesSkipsMalformatIntEnvVariable(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	expected := 123
	c.Server.RpcPort = expected
	os.Setenv("LFSSYNC_SERVER_RPC_PORT", "not_an_int")
	c.OverrideFromEnvironment()
	if c.Server.RpcPort != expected {
		t.Fatalf(
			`%s:, Expected Server.RpcPort = %d, found: %d`,
			issue, expected, c.Server.RpcPort)
	}
}

func TestOverridesBoolEnvVariable(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	expected := true
	os.Setenv("LFSSYNC_DATABASE_MIGRATE", "true")
	c.OverrideFromEnvironment()
	if c.Database.Migrate != expected {
		t.Fatalf(
			`%s:, Expected Database.Migrate = %t, found: %t`,
			issue, expected, c.Database.Migrate)
	}
}
