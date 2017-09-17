package RESTHandlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store/testing"
)

func init() {
	logger.Log.Out = ioutil.Discard // Logs redirected to /dev/null
}

func TestHandler_GetObjectMockedStore(t *testing.T) {
	t.Parallel()
	storage := &storeTesting.StorageMock{}

	r := ApiHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name         string
		storageValue *store.Object
		storageErr   error
		responseCode int
		responseBody []byte
	}{
		{"KeyNotFound", nil, store.KeyNotFound, http.StatusNotFound, []byte{}},
		{"StorageIsLocked", nil, store.StorageIsLocked, http.StatusInternalServerError, []byte{}},
		{"ReceivedStop", nil, store.ReceivedStop, http.StatusInternalServerError, []byte{}},
		{"Success", &store.Object{"key", "MediaType", []byte("foo")}, nil, http.StatusOK, []byte("foo")},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			storage.ErrOnGet = c.storageErr
			storage.ResponseOnGet = c.storageValue

			resp, err := http.Get(srv.URL + "/objects/adcae")
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

			if resp.StatusCode == http.StatusOK {
				ct, ok := resp.Header["Content-Type"]
				if !ok || len(ct) == 0 || resp.Header["Content-Type"][0] != c.storageValue.MediaType {
					t.Error("Content-Type was corrupted")
				}
			}
		})
	}
}

func TestHandler_GetObjectRealStore(t *testing.T) {
	t.Parallel()
	expectedValue := []byte("A long time ago, in a galaxy far, far away…")
	expectedMediaType := "MediaType"

	storage := store.NewStorage()
	id, err := storage.Create(context.Background(), expectedValue, expectedMediaType)
	if err != nil {
		t.Fatal(err)
	}

	r := ApiHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/objects/" + id)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d status", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(body, expectedValue) {
		t.Error("Response data was corrupted")
	}

	ct, ok := resp.Header["Content-Type"]
	if !ok || len(ct) == 0 || resp.Header["Content-Type"][0] != expectedMediaType {
		t.Error("Content-Type was corrupted")
	}
}
