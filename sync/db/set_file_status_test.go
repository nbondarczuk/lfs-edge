package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// We expect one row with changed status - positive test.
func TestSetStatus(t *testing.T) {
	// Build test case config
	var (
		testFileId    uint64 = 2
		testFileName  string = "2"
		testStatus    int32  = 123
		testNewStatus int32  = 1234
	)
	err := CreateFileInfo(testFileId, testFileName, testStatus)
	assert.Nil(t, err)

	// Make test
	err = SetFileStatus(testFileId, testNewStatus)

	// Check results
	assert.Nil(t, err)
	entry, err := GetFile(testFileId)
	assert.Nilf(t, err, "GetFile() with no error, got: %v", err)
	assert.Truef(t, entry.Exist, "GetFile() with id: %d, expected one entry",
		testFileId)
	assert.Equalf(t, entry.Name, testFileName,
		"GetFile() with id: %d, expected file name: %s, got: %s",
		testFileId, testFileName, entry.Name)
	assert.Equalf(t, entry.Status, testNewStatus,
		"GetFile() with id: %d, expected status: %d, got: %d",
		testFileId, testNewStatus, entry.Status)
}

// We expect one row with changed status and metadata - positive test.
func TestSetStatusWithMetadata(t *testing.T) {
	// Build test case config
	var (
		testFileId    uint64 = 3
		testFileName  string = "3"
		testStatus    int32  = 1234
		testNewStatus int32  = 12345
		testNewName   string = "xxx"
		testNewSize   int64  = 666
	)
	err := CreateFileInfo(testFileId, testFileName, testStatus)
	assert.Nil(t, err)

	// Make test
	err = SetFileStatusWithMetadata(testFileId, testNewStatus, testNewName, 666, "whatever", time.Now())

	// Check results
	assert.Nil(t, err)
	entry, err := GetFile(testFileId)
	assert.Nilf(t, err, "GetFile() with no error, got: %v", err)
	assert.Truef(t, entry.Exist, "GetFile() with id: %d, expected one entry",
		testFileId)
	assert.Equalf(t, entry.Status, testNewStatus,
		"SetFileStatusWithMetadata() with id: %d, expected status: %d, got: %d",
		testFileId, testNewStatus, entry.Status)
	assert.Equalf(t, entry.Name, testNewName,
		"SetFileStatusWithMetadata() with id: %d, expected name: %s, got: %s",
		testFileId, testNewName, entry.Name)
	assert.Equalf(t, entry.Size, testNewSize,
		"SetFileStatusWithMetadata() with id: %d, expected size: %d, got: %d",
		testFileId, testNewSize, entry.Size)	
}
