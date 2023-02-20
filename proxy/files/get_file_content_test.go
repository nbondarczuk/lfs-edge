package files_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"
)

// The tests document the semantics of the file content access function
func TestGetFileContent(t *testing.T) {
	mockConfig(t, "/tmp")

	t.Run("errors get contents of missing file with incorrect id", func(t *testing.T) {
		_, err := files.GetFileContent("request", "device", "whatever")
		assert.NotNil(t, err)
		assert.Equal(t, err, files.ErrInvalidArg)
	})

	t.Run("gets contents of the existing correct file", func(t *testing.T) {
		// Prepare test data
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/0", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Check the contents of the file
		d1, err := os.ReadFile("/tmp/0")
		if err != nil {
			panic(err)
		}
		assert.Equal(t, string(d1), msg)

		// Make the test and check the result
		fc, err := files.GetFileContent("request", "device", "0")
		assert.Nil(t, err)
		assert.Equal(t, fc.FileName, "/tmp/0")
		assert.NotNil(t, fc.ServedFile)
		defer fc.ServedFile.Close()

		// Check the contents using file descriptor used for reading
		d2 := make([]byte, len(msg))
		n, err := fc.ServedFile.Read(d2)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, n, len(msg))
		assert.Equal(t, fmt.Sprintf(string(d2)), msg)

		// Cleanup
		err = os.Remove("/tmp/0")
		if err != nil {
			panic(err)
		}
	})

	t.Run("errors get contents of missing file with correct id", func(t *testing.T) {
		_, err := files.GetFileContent("request", "device", "0")
		assert.NotNil(t, err)
	})	
}
