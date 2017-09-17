package store

import "testing"

func TestStoreErr_Error(t *testing.T) {
	if ValueSizeIsExceeded.Error() != string(ValueSizeIsExceeded) {
		t.Error("Expected equal values")
	}
}
