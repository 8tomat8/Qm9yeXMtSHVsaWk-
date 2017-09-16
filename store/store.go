package store

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

const MaxValueSize = 1 << 20

type Storage interface {
	GetAll() []*Object
	Get(string) (*Object, error)
	Create([]byte, string) (string, error)
	Delete(string) error
}

type storage struct {
	// TODO: Rewrite this piece of... map...
	data map[string]*Object
	muRW sync.RWMutex
}

func GetStorage() Storage {
	once.Do(func() {
		mockedStore = storage{
			data: make(map[string]*Object),
			muRW: sync.RWMutex{},
		}
	})
	return Storage(&mockedStore)
}

var (
	mockedStore storage
	once        sync.Once
)

type Object struct {
	Key       string
	MediaType string
	Value     []byte
}

func (s *storage) GetAll() []*Object {
	s.muRW.RLock()
	defer s.muRW.RUnlock()

	result := make([]*Object, 0, len(s.data))

	for i := range s.data {
		result = append(result, s.data[i])
	}
	return result
}

func (s *storage) Get(key string) (*Object, error) {
	s.muRW.RLock()
	defer s.muRW.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return nil, KeyNotFound
	}

	return value, nil
}

func (s *storage) Create(value []byte, mediaType string) (string, error) {
	h := md5.New()

	n, err := h.Write(value)
	if err != nil {
		return "", err
	} else if n != len(value) {
		return "", StoreErr("couldn't count hash for value")
	}

	hashSum := hex.EncodeToString(h.Sum(nil))

	s.muRW.Lock()
	defer s.muRW.Unlock()

	_, ok := s.data[hashSum]
	if ok {
		return "", KeyAlreadyExist
	}

	s.data[hashSum] = &Object{Key: hashSum, MediaType: mediaType, Value: value}
	return hashSum, nil
}

func (s *storage) Delete(key string) error {
	s.muRW.Lock()
	defer s.muRW.Unlock()

	// TODO: Resolve "read before write"
	_, ok := s.data[key]
	if !ok {
		return KeyNotFound
	}
	delete(s.data, key)
	return nil
}
