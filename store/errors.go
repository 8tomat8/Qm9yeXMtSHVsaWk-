package store

type StoreErr string

const (
	KeyAlreadyExist = StoreErr("Key with such name is already exist")
	KeyNotFound     = StoreErr("Key with such name was not found in store")
)

func (err StoreErr) Error() string {
	return string(err)
}
