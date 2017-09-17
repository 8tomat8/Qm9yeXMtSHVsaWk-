package store

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

const (
	// MaxValueSize stores maximum allowed size for value
	MaxValueSize            = 1 << 20
	defaultStoreLockTimeout = 5
)

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

type Object struct {
	Key       string
	MediaType string
	Value     []byte
}

// Get will return value from storage by its key
// If key does not exist KeyNotFound error sill be returned
// Function handles ctx.Done and could be canceled before finish. In this case ReceivedStop error will be returned
func (s *storage) Get(ctx context.Context, key string) (*Object, error) {
	unlock, err := s.rlock(ctx)
	if err != nil {
		return nil, err
	}
	defer unlock()

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
	// Also md5 is quite complex algorithm, but due to task requirements could not be changed to some simpler
	h := md5.New()

	n, err := h.Write(value)
	if err != nil {
		return "", err
	} else if n != len(value) {
		return "", StoreErr("couldn't count hash for value")
	}

	hashSum := hex.EncodeToString(h.Sum(nil))

	unlock, err := s.lock(ctx)
	if err != nil {
		return "", err
	}

	_, ok := s.data[hashSum]
	if ok {
		unlock()
		return "", KeyAlreadyExist
	}

	s.data[hashSum] = &Object{Key: hashSum, MediaType: mediaType, Value: value}
	unlock()
	return hashSum, nil
}

// Delete function implements logic to remove object from storage by its key.
// If key does not exist KeyNotFound error sill be returned
// Function handles ctx.Done and could be canceled before finish. In this case ReceivedStop error will be returned
func (s *storage) Delete(ctx context.Context, key string) error {
	unlock, err := s.lock(ctx)
	if err != nil {
		return err
	}
	unlock()

	// TODO: Resolve "read before write"
	_, ok := s.data[key]
	if !ok {
		return KeyNotFound
	}
	delete(s.data, key)
	return nil
}

func (s *storage) lock(ctx context.Context) (func(), error) {
	// Extra code to handle deadlocks and ctx.Done
	locked := make(chan struct{})
	go func() {
		s.muRW.Lock()
		locked <- empty
	}()

	ticker := time.NewTicker(time.Second * defaultStoreLockTimeout)
	defer ticker.Stop()
	select {
	case <-locked:
	case <-ticker.C:
		return s.muRW.Unlock, StorageIsLocked
	case <-ctx.Done():
		return s.muRW.Unlock, ReceivedStop
	}
	return s.muRW.Unlock, nil
}

func (s *storage) rlock(ctx context.Context) (func(), error) {
	// Extra code to handle deadlocks and ctx.Done
	locked := make(chan struct{})
	go func() {
		s.muRW.RLock()
		locked <- empty
	}()

	ticker := time.NewTicker(time.Second * defaultStoreLockTimeout)
	defer ticker.Stop()
	select {
	case <-locked:
	case <-ticker.C:
		return s.muRW.RUnlock, StorageIsLocked
	case <-ctx.Done():
		return s.muRW.RUnlock, ReceivedStop
	}
	return s.muRW.RUnlock, nil
}
