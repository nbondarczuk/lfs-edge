package files

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common"	
)

const (
	newFileAccessRights = 0644
)

func CheckFileExistWithUrl(uri string) (bool, error) {
	// Make HEAD request to check if file exists
    r, err := http.Head(uri)
    if err != nil {
        return false, err
    }

	var ok bool // default is false
	if r.StatusCode == http.StatusOK {
		ok = true
	}
	
	return ok, nil
}

// downloadPayloadWithUrl issues a GET request to the server.
func downloadPayloadWithUrl(uri string) ([]byte, error) {
	// Make GET request for file.
    r, err := http.Get(uri)
    if err != nil {
        return nil, err
    }
    defer r.Body.Close()
	
	// Read payload from request.
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil, err
    }

    return data, nil
}

func savePayload(path string, data []byte) error {
    return ioutil.WriteFile(path, data, newFileAccessRights)
}

// DownloadFile downloads file from from server with GET request
// and saves it to given path.
func DownloadFileWithUrl(uri string, path string) (err error) {
	data, err := downloadPayloadWithUrl(uri)
	if err != nil {
		return err
	}

	err = savePayload(path, data)
	if err != nil {
		return err
	}
	
	return nil
}

// DownloadFileWithPresignedUrl first gets the presigned url from FS and then
// it uses it to download file content.
func DownloadFileWithPresignedUrl(template string, id uint64, path string) error {
	presignedUrl := fmt.Sprintf(template, id)

	payload, err := downloadPayloadWithUrl(presignedUrl)
	if err != nil {
		return err
	}

	var response common.SignedUrlResponse
	err = json.Unmarshal(payload, &response)
	if err != nil {
		return err
	}
	uri := response.SignedUrl

	ok, err := CheckFileExistWithUrl(uri)
	if err != nil {
		return err
	}	
	if !ok {
		return ErrNotFound
	}

	return DownloadFileWithUrl(uri, path)
}
