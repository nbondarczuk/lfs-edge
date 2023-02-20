package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// We dont expect any rows with id 0 - negative test.
func TestGetFileWithNonExistentId(t *testing.T) {
	// Make test
	entry, err := GetFile(0)

	// check results
	assert.Nilf(t, err, "GetFile() with 0 id: found: %v", err)
	assert.Falsef(t, entry.Exist, "GetFile() with 0 id: not expected")
}

// We expect one row with id 1 - positive test.
func TestGetFileWithExistentId(t *testing.T) {
	// Build test case config
	var (
		testFileId   uint64 = 1
		testFileName string = "xxx"
		testStatus   int32  = 12
	)
	err := CreateFileInfo(testFileId, testFileName, testStatus)
	assert.Nil(t, err)

	// Make test
	entry, err := GetFile(testFileId)

	// Check results
	assert.Nilf(t, err, "GetFile() with no error, got: %v", err)
	assert.Truef(t, entry.Exist, "GetFile() with id: %d, expected one entry",
		testFileId)
	assert.Equalf(t, entry.Name, testFileName,
		"GetFile() with id: %d, expected file name: %s, got: %s",
		testFileId, testFileName, entry.Name)
	assert.Equalf(t, entry.Status, testStatus,
		"GetFile() with id: %d, expected status: %d, got: %d",
		testFileId, testStatus, entry.Status)

}
