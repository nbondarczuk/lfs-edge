package file

import (
	"crypto/md5"	
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	WholeFileContent = -1
	DirPath          = 1
	RegularPath      = 2
)

// PathExists checks if the file dir exists and it is a dir using stat so
// no need to close it.
func PathExists(path string, typeOfPath uint) (bool, error) {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return true, err
	}

	// Interpret the file type flag, it may be dir or reguular file
	if typeOfPath == DirPath {
		if !fi.Mode().IsDir() {
			return false, nil
		}
	} else if typeOfPath == RegularPath {
		if fi.Mode().IsDir() {
			return false, nil
		}
	} else {
		return false, ErrInvalidArg
	}

	// Successful case, the file path exists and is as specified, dir or regular
	return true, nil
}

// StaticInfo gets stat of the file, created as it seems not in all Unix brands.
// Due to the fact of multiple return values named return makes it more readable
// as there is no need to explicitely return zero values or create empty structures
// just to handle errors.
func StaticInfo(path string) (size int64, modified time.Time, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	size = fi.Size()
	modified = fi.ModTime()

	return
}

// fileChecksumSHA256 gets SHA256 of the file using open so close must follow
func ChecksumSHA256Info(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	cksum := fmt.Sprintf("%x", h.Sum(nil))

	return cksum, nil
}

// ChecksumMD5Info gets MD5 of the file using open so close must follow
func ChecksumMD5Info(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	cksum := fmt.Sprintf("%x", h.Sum(nil))

	return cksum, nil
}
