package files_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"
)

// The tests document the semantics of the file info function
func TestGetFileInfo(t *testing.T) {
	mockConfig(t, "/tmp")

	t.Run("errors get info of missing file with incorrect id", func(t *testing.T) {
		_, err := files.GetFileInfo("xxx")
		assert.NotNil(t, err)
		assert.Equal(t, err, files.ErrInvalidArg)
	})

	t.Run("errors get contents of missing file with correct id", func(t *testing.T) {
		_, err := files.GetFileInfo("0")
		assert.NotNil(t, err)
		assert.Equal(t, err, files.ErrNotFound)
	})

	t.Run("gets info of the existing correct file", func(t *testing.T) {
		// Prepare test data
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/1", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Check the contents of the file
		d1, err := os.ReadFile("/tmp/1")
		if err != nil {
			panic(err)
		}
		assert.Equal(t, string(d1), msg)

		// Make the test and check the result
		fc, err := files.GetFileInfo("1")
		assert.Nil(t, err)
		assert.Equal(t, fc.FileID, uint64(1))
		assert.NotEqual(t, fc.Checksum, "")
		assert.Equal(t, fc.Size, int64(len(msg)))
		assert.NotEqual(t, fc.UpdatedAt, time.Time{})

		// Cleanup
		err = os.Remove("/tmp/1")
		if err != nil {
			panic(err)
		}
	})
}
