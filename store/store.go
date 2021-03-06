package store

import (
	"context"
	"encoding/hex"
	"sync"

	"github.com/spaolacci/murmur3"
)

// MaxValueSize stores maximum allowed size for value
const MaxValueSize = 1 << 20

var empty struct{}

// Storage describes storage interface
type Storage interface {
	Get(ctx context.Context, key string) (*Object, error)
	Create(ctx context.Context, value []byte, mediaType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type storage struct {
	// TODO: Rewrite this piece of... map...
	data map[string]*Object
	muRW sync.RWMutex
}

// NewStorage will initialize new empty storage and return it
func NewStorage() Storage {
	return Storage(&storage{
		data: make(map[string]*Object),
		muRW: sync.RWMutex{},
	})
}

// Object represents structure of each object in store
type Object struct {
	Key       string
	MediaType string
	Value     []byte
}

// Get will return value from storage by its key
// If key does not exist KeyNotFound error sill be returned
// Function handles ctx.Done and could be canceled before finish. In this case ReceivedStop error will be returned
func (s *storage) Get(ctx context.Context, key string) (*Object, error) {
	err := s.lock(ctx, s.muRW.RLock)
	if err != nil {
		return nil, err
	}
	defer s.muRW.RUnlock()

	object, ok := s.data[key]
	if !ok {
		return nil, KeyNotFound
	}

	return object, nil
}

// Create function implements logic to save data in storage.
// If value is not unique it will return KeyAlreadyExist error
// Function handles ctx.Done and could be canceled before finish. In this case ReceivedStop error will be returned
func (s *storage) Create(ctx context.Context, value []byte, mediaType string) (string, error) {
	if len(value) > MaxValueSize {
		return "", ValueSizeIsExceeded
	}

	// Hash is prone to collision. Fast but not optimal solution
	h := murmur3.New128()

	n, err := h.Write(value)
	if err != nil {
		return "", err
	} else if n != len(value) {
		return "", Err("couldn't count hash for value")
	}

	hashSum := hex.EncodeToString(h.Sum(nil))

	err = s.lock(ctx, s.muRW.Lock)
	if err != nil {
		return "", err
	}

	_, ok := s.data[hashSum]
	if ok {
		s.muRW.Unlock()
		return "", KeyAlreadyExist
	}

	s.data[hashSum] = &Object{Key: hashSum, MediaType: mediaType, Value: value}
	s.muRW.Unlock()
	return hashSum, nil
}

// Delete function implements logic to remove object from storage by its key.
// If key does not exist KeyNotFound error sill be returned
// Function handles ctx.Done and could be canceled before finish. In this case ReceivedStop error will be returned
func (s *storage) Delete(ctx context.Context, key string) error {
	err := s.lock(ctx, s.muRW.Lock)
	if err != nil {
		return err
	}
	defer s.muRW.Unlock()

	// TODO: Resolve "read before write"
	_, ok := s.data[key]
	if !ok {
		return KeyNotFound
	}
	delete(s.data, key)
	return nil
}

// lock is a helper function that allows to handle lock
func (s *storage) lock(ctx context.Context, locker func()) error {
	// Extra code to handle deadlocks and ctx.Done
	locked := make(chan struct{})
	go func() {
		locker()
		locked <- empty
	}()

	select {
	case <-locked:
	case <-ctx.Done():
		return ReceivedStop
	}
	return nil
}
