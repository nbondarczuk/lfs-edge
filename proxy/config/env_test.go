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
	c := Config{}
	expected := c.Server.Host
	c.OverrideFromEnvironment()
	if c.Server.Host != "" {
		t.Fatalf(
			`OverrideFromEnvironment: Expected Server.Host = %s,  found: %s`,
			expected, c.Server.Host)
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
			"LFSPROXY_HOST", "Server.Host",
			&c.Server.Host, "localhost", "localhost",
		},
		{
			"LFSPROXY_RPC_HOST", "RpcServer.Host",
			&c.RpcServer.Host, "whatever-host", "whatever-host",
		},
		{
			"LFSPROXY_FILES_PATH", "Files.Path",
			&c.Files.Path, "whatever-path", "whatever-path",
		},
	}

	for _, tc := range tests {
		os.Setenv(tc.envVarName, tc.inputValue)
		c.OverrideFromEnvironment()
		if *tc.configVar != tc.expectedValue {
			t.Fatalf(
				`%s: Expected %s = %s, found: %s`,
				issue, tc.fieldVarName, tc.expectedValue, *tc.configVar,
			)
		}
	}
}

func TestOverridesAllIntEnvVariables(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	type test struct {
		envVarName    string
		fieldVarName  string
		configVar     *int
		inputValue    int64
		expectedValue int
	}
	tests := []test{
		{
			"LFSPROXY_PORT", "Server.Port",
			&c.Server.Port, 123, 123,
		},
		{
			"LFSPROXY_RPC_PORT", "RpcServer.Port",
			&c.RpcServer.Port, 1234, 1234,
		},
	}

	for _, tc := range tests {
		os.Setenv(tc.envVarName, strconv.FormatInt(tc.inputValue, 10))
		c.OverrideFromEnvironment()
		if *tc.configVar != tc.expectedValue {
			t.Fatalf(
				`%s: Expected %s = %d, found: %d`,
				issue, tc.fieldVarName, tc.expectedValue, *tc.configVar,
			)
		}
	}
}

func TestOverridesSkipsBadEnvVariable(t *testing.T) {
	issue := "OverrideFromEnvironment"
	c := Config{}
	expected := 123
	c.Server.Port = expected
	os.Setenv("LFSPROXY_PORT", "not_an_int")
	c.OverrideFromEnvironment()
	if c.Server.Port != expected {
		t.Fatalf(
			`%s: Expected Server.Port = %d, found: %d`,
			issue, expected, c.Server.Port,
		)
	}
}
