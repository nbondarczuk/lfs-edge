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
func TestGetFileContent(t *testing.T) {
	mockConfig(t, "/tmp")

	t.Run("gets content of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/5", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/5?device=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read whole paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, len(data))
		assert.Equal(t, msg, string(data))

		// Cleanup
		err = os.Remove("/tmp/5")
		if err != nil {
			panic(err)
		}
	})

	t.Run("gets range content of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/6", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/6?device=abc", nil)
		req.Header.Set("Range", "bytes=0-9")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read whole paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, len(data))
		assert.Equal(t, msg, string(data))

		// Cleanup
		err = os.Remove("/tmp/6")
		if err != nil {
			panic(err)
		}
	})


	t.Run("gets prefix range of content of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/7", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/7?device=abc", nil)
		req.Header.Set("Range", "bytes=0-0")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, len(data), 0)
		assert.Equal(t, len(data), 1)
		assert.Equal(t, "m", string(data))

		// Cleanup
		err = os.Remove("/tmp/7")
		if err != nil {
			panic(err)
		}
	})

	t.Run("gets suffix range of content of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/8", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/8?device=abc", nil)
		req.Header.Set("Range", "bytes=9-9")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, len(data), 0)
		assert.Equal(t, len(data), 1)
		assert.Equal(t, "d", string(data))

		// Cleanup
		err = os.Remove("/tmp/8")
		if err != nil {
			panic(err)
		}
	})

	t.Run("gets infix range of content of an existing file", func(t *testing.T) {
		var msg string = "melloworld"
		err := os.WriteFile("/tmp/9", []byte(msg), 0644)
		if err != nil {
			panic(err)
		}

		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/9?device=abc", nil)
		req.Header.Set("Range", "bytes=5-6")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Read results of the request
		result := w.Result()
		defer result.Body.Close()

		// Read paylad of the reply and check it for success
		data, err := ioutil.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.NotEqual(t, len(data), 0)
		assert.Equal(t, len(data), 2)
		assert.Equal(t, "wo", string(data))

		// Cleanup
		err = os.Remove("/tmp/9")
		if err != nil {
			panic(err)
		}
	})

	t.Run("errors get content of a missing file", func(t *testing.T) {
		// Make router for the handler being tested
		r := mux.NewRouter()
		r.HandleFunc("/api/v1/files/{id:[0-9]+}", rest.GetFileContentHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/10?device=abc", nil)
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
}
