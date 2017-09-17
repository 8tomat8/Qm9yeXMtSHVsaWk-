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

func TestHandler_DeleteObjectMockedStore(t *testing.T) {
	t.Parallel()
	logger.Log.Out = ioutil.Discard // Logs redirected to /dev/null

	storage := &storeTesting.StorageMock{}

	r := ApiHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	HTTPclient := &http.Client{}

	testCases := []struct {
		name         string
		storageErr   error
		responseCode int
		responseBody []byte
	}{
		{"KeyNotFound", store.KeyNotFound, http.StatusNotFound, nil},
		{"ReceivedStop", store.ReceivedStop, http.StatusInternalServerError, nil},
		{"Success", nil, http.StatusNoContent, nil},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			storage.ErrOnDelete = c.storageErr

			req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/objects/adcae", nil)
			resp, err := HTTPclient.Do(req)
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

func TestHandler_DeleteObjectRealStore(t *testing.T) {
	t.Parallel()
	logger.Log.Out = ioutil.Discard // Logs redirected to /dev/null

	expectedMediaType := "MediaType"

	storage := store.NewStorage()
	id, err := storage.Create(context.Background(), []byte("foo"), expectedMediaType)
	if err != nil {
		t.Fatal(err)
	}

	r := ApiHandler(storage)
	srv := httptest.NewServer(r)
	defer srv.Close()

	HTTPclient := &http.Client{}

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/objects/"+id, nil)
	resp, err := HTTPclient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected %d, got %d status", http.StatusNoContent, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body) != 0 {
		t.Error("Expected empty response")
	}
}
