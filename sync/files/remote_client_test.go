package files_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/files"
)

type client struct {
	url string
}

func newClient(url string) client {
	return client{url}
}

func (c client) tryCheckExist(path string) (bool, error) {
	return files.CheckFileExistWithUrl(c.url + path)
}

// Tests function checking HTTP server file existence.
func TestClientCheckExist(t *testing.T) {
	// Test data
	var (
		testPathName    = "/tmp/1"
		testFilePayload = "hello"
	)

	// Prepare test env.
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testFilePayload)
	}))
	defer svr.Close()
	c := newClient(svr.URL)
	t.Cleanup(func() {os.Remove(testPathName)})
	
	// Run test on client.
	ok, err := c.tryCheckExist(testPathName)

	// Check test results.
	assert.Nil(t, err)
	assert.Truef(t, ok, 
		"expected result of exist check to be true for: %s%s",
		svr.URL, testPathName)
}

func (c client) tryDownloadFile(id uint64, path string) error {
	return files.DownloadFileWithUrl(c.url, path)
}

// Tests function downloading HTTP server file content.
func TestClientDownloadFile(t *testing.T) {
	// Test data
	var (
		testFileID      uint64 = 2
		testFileName           = "/tmp/2"
		testFilePayload        = "hello"
	)

	// Prepare test env.
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testFilePayload)
	}))
	defer svr.Close()
	c := newClient(svr.URL)
	t.Cleanup(func() {os.Remove(testFileName)})
	
	// Run test on client.
	err := c.tryDownloadFile(testFileID, testFileName)

	// Check test results.
	assert.Nil(t, err)
	assert.FileExists(t, testFileName)
	content, err := os.ReadFile(testFileName)
	assert.Nil(t, err)
	assert.Equalf(t, string(content), testFilePayload,
		"invalid result file: %s content: %s, expected: %s",
		testFileName, content, testFilePayload)
}

