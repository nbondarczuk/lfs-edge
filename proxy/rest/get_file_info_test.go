package rest_test

import (
	//"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rest"
)

// The tests document the semantics of the file content access function
func TestGetFileInfo(t *testing.T) {
	mockConfig(t, "/tmp")

	t.Run("errors get info of a missing file", func(t *testing.T) {
		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileInfoHandler)

		// Make request
		req := httptest.NewRequest(http.MethodHead, "/api/v1/files/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read paylad of the reply and check it for error
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, len(data), 0)
	})

	t.Run("gets info of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/3", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileInfoHandler)

		// Make request
		req := httptest.NewRequest(http.MethodHead, "/api/v1/files/3", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, len(data), 0)

		// Cleanup
		err = os.Remove("/tmp/3")
		if err != nil {
			panic(err)
		}
	})
}
