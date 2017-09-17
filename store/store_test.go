package store

import (
	"bytes"
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	s := NewStorage()
	ctx := context.Background()

	testValue := []byte("A wizard is never late, Frodo Baggins.")
	testMediaType := "foo bar baz"

	id, err := s.Create(ctx, testValue, testMediaType)
	if err != nil {
		t.Error("Received error on Create:\n", err)
	}

	object, err := s.Get(ctx, id)
	if err != nil {
		t.Error("Received error on Get:\n", err)
	}

	if !bytes.Equal(object.Value, testValue) {
		t.Error("Stored value was corrupted")
	}

	if object.MediaType != testMediaType {
		t.Error("Media type was corrupted")
	}

	err = s.Delete(ctx, id)
	if err != nil {
		t.Error("Received error on Delete:\n", err)
	}

	_, err = s.Get(ctx, id)
	if err != KeyNotFound {
		t.Error("Delete did not delete obect")
	}
}

func TestLock(t *testing.T) {
	s := NewStorage().(*storage)
	ctx, cancel := context.WithCancel(context.Background())

	s.muRW.Lock()
	cancel()
	err := s.lock(ctx, s.muRW.Lock)
	if err != ReceivedStop {
		t.Error("lock function does not handle long lock properly")
	}
	s.muRW.Unlock()
}

func TestStorage_CreateMaxValueSize(t *testing.T) {
	s := NewStorage()

	bigValue := make([]byte, MaxValueSize+1)
	rand.Read(bigValue)

	_, err := s.Create(context.Background(), bigValue, "foo")
	if err != ValueSizeIsExceeded {
		t.Error("Create allowed to create object with size more than MaxValueSize")
	}
}

func TestStorage_CreateKeyAlreadyExist(t *testing.T) {
	s := NewStorage()
	rand.Seed(time.Now().UnixNano())

	value := make([]byte, 1000)

	for i := 0; i < 1000; i++ {
		rand.Read(value)
		_, err := s.Create(context.Background(), value, "foo")
		if err != nil {
			t.Error("Expected nil error")
		}

		_, err = s.Create(context.Background(), value, "foo")
		if err != KeyAlreadyExist {
			t.Error("Expected KeyAlreadyExist error")
		}
	}

}

func TestStorage_CreateHashCollision(t *testing.T) {
	s := NewStorage()
	rand.Seed(time.Now().UnixNano())

	value := make([]byte, 10000)
	for i := 0; i < 1000; i++ {
		rand.Read(value)
		_, err := s.Create(context.Background(), value, "foo")
		if err != nil {
			t.Error("Expected nil error")
		}

		rand.Read(value)
		_, err = s.Create(context.Background(), value, "foo")
		if err == KeyAlreadyExist {
			t.Error("Collision detected. Hash algorithm should be changed")
		}
	}

}
