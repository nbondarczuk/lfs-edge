package file_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
)

// The tests document the semantic of the file access functions
func TestFile(t *testing.T) {
	t.Run("errors check on missing dir path", func(t *testing.T) {
		ok, err := file.PathExists("/whatever", file.DirPath)
		assert.Nil(t, err)
		assert.Equal(t, ok, false)
	})

	t.Run("errors check on missing regular path", func(t *testing.T) {
		ok, err := file.PathExists("/whatever", file.RegularPath)
		assert.Nil(t, err)
		assert.Equal(t, ok, false)
	})

	t.Run("detects existing dir folder", func(t *testing.T) {
		ok, err := file.PathExists("/tmp", file.DirPath)
		assert.Nil(t, err)
		assert.Equal(t, ok, true)
	})

	t.Run("detects existing regular file", func(t *testing.T) {
		ok, err := file.PathExists("/etc/hosts", file.RegularPath)
		assert.Nil(t, err)
		assert.Equal(t, ok, true)
	})

	t.Run("errors stat on missing file", func(t *testing.T) {
		_, _, err := file.StaticInfo("/whatever")
		assert.NotNil(t, err, nil)
	})

	t.Run("gets stat on regular file", func(t *testing.T) {
		size, modified, err := file.StaticInfo("/etc/hosts")
		assert.Nil(t, err)
		assert.NotEqual(t, modified, time.Time{})
		assert.NotEqual(t, size, 0)
	})

	t.Run("errors cksum on missing file", func(t *testing.T) {
		_, err := file.ChecksumSHA256Info("/whatever")
		assert.NotNil(t, err)
	})

	t.Run("gets cksum on regular file", func(t *testing.T) {
		cksum, err := file.ChecksumSHA256Info("/etc/hosts")
		assert.Nil(t, err)
		assert.NotEqual(t, len(cksum), 0)
	})
}
