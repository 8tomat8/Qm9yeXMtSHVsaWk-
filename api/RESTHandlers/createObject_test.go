package RESTHandlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"math/rand"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store/testing"
)

func init() {
	logger.Log.Out = ioutil.Discard // Logs redirected to /dev/null
}

func TestHandler_CreateObjectMockedStore(t *testing.T) {
	t.Parallel()
	storage := &storeTesting.StorageMock{}

	r := APIHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	bigValue := make([]byte, store.MaxValueSize+1)
	rand.Read(bigValue)

	testCases := []struct {
		name         string
		storageKey   string
		storageErr   error
		requestBody  []byte
		responseCode int
		responseBody []byte
	}{
		{"KeyAlreadyExist", "", store.KeyAlreadyExist, []byte("foo"), http.StatusConflict, []byte{}},
		{"ValueSizeIsExceeded", "", store.ValueSizeIsExceeded, []byte("foo"), http.StatusRequestEntityTooLarge, []byte{}},
		{"ValueSizeIsExceeded", "", nil, bigValue, http.StatusRequestEntityTooLarge, []byte{}},
		{"ReceivedStop", "", store.ReceivedStop, []byte("foo"), http.StatusInternalServerError, []byte{}},
		{"Success", "keyName", nil, []byte("foo"), http.StatusCreated, []byte("keyName")},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			storage.ErrOnCreate = c.storageErr
			storage.ResponseOnCreate = c.storageKey

			resp, err := http.Post(srv.URL+"/objects", "text/plain", bytes.NewReader(c.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != c.responseCode {
				t.Errorf("Expected %d, got %d status", c.responseCode, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(body, c.responseBody) {
				t.Error("Response data was corrupted")
			}
		})
	}
}

func TestHandler_CreateObjectRealStore(t *testing.T) {
	t.Parallel()
	expectedValue := []byte("A long time ago, in a galaxy far, far away…")
	expectedMediaType := "plan/text"

	storage := store.NewStorage()
	r := APIHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/objects", expectedMediaType, bytes.NewReader(expectedValue))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected %d, got %d status", http.StatusCreated, resp.StatusCode)
	}

	id, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	resultObject, err := storage.Get(context.Background(), string(id))
	if err != nil {
		t.Error("Object was not stored. Storage returned error:\n", err)
	}

	if !bytes.Equal(resultObject.Value, expectedValue) {
		t.Error("Stored value was corrupted")
	}

	if resultObject.MediaType != expectedMediaType {
		t.Error("Media type was corrupted")
	}

}

func TestHandler_CreateObjectNoContentType(t *testing.T) {
	t.Parallel()
	storage := &storeTesting.StorageMock{}

	r := APIHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	bigValue := make([]byte, store.MaxValueSize+1)
	rand.Read(bigValue)

	resp, err := http.Post(srv.URL+"/objects", "", bytes.NewReader([]byte("foo")))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Expected %d, got %d status", http.StatusUnsupportedMediaType, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body) != 0 {
		t.Error("Expected empty response body")
	}
}
