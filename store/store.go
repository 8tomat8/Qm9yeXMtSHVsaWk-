package store

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

const (
	MaxValueSize            = 1 << 20
	defaultStoreLockTimeout = 5
)

var empty struct{}

type Storage interface {
	Get(context.Context, string) (*Object, error)
	Create(context.Context, []byte, string) (string, error)
	Delete(context.Context, string) error
}

type storage struct {
	// TODO: Rewrite this piece of... map...
	data    map[string]*Object
	hashSet map[uint64]struct{}
	muRW    sync.RWMutex
}

func NewStorage() Storage {
	return Storage(&storage{
		data:    make(map[string]*Object),
		hashSet: make(map[uint64]struct{}),
		muRW:    sync.RWMutex{},
	})
}

type Object struct {
	Key       string
	MediaType string
	Value     []byte
}

func (s *storage) Get(ctx context.Context, key string) (*Object, error) {
	unlock, err := s.rlock(ctx)
	if err != nil {
		return nil, err
	}

	value, ok := s.data[key]
	if !ok {
		unlock()
		return nil, KeyNotFound
	}

	unlock()
	return value, nil
}

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

func (s *storage) Delete(ctx context.Context, key string) error {
	unlock, err := s.lock(ctx)
	if err != nil {
		return err
	}

	// TODO: Resolve "read before write"
	_, ok := s.data[key]
	if !ok {
		return KeyNotFound
	}
	delete(s.data, key)
	unlock()
	return nil
}

func (s *storage) lock(ctx context.Context) (func(), error) {
	// Extra code to handle deadlock
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
	// Extra code to handle deadlock
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
