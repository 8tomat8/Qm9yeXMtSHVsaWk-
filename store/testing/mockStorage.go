package storeTesting

import (
	"context"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
)

// StorageMock is the simplest implementation of Storage interface
type StorageMock struct {
	ResponseOnCreate string
	ErrOnCreate      error

	ResponseOnGet *store.Object
	ErrOnGet      error

	ErrOnDelete error
}

// Get will return ResponseOnGet and ErrOnGet from StorageMock
func (s *StorageMock) Get(context.Context, string) (*store.Object, error) {
	return s.ResponseOnGet, s.ErrOnGet
}

// Create will return ResponseOnCreate and ErrOnCreate from StorageMock
func (s *StorageMock) Create(context.Context, []byte, string) (string, error) {
	return s.ResponseOnCreate, s.ErrOnCreate
}

// Delete will return ErrOnDelete from StorageMock
func (s *StorageMock) Delete(context.Context, string) error {
	return s.ErrOnDelete
}
