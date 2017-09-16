package store

type StoreErr string

const (
	KeyAlreadyExist     = StoreErr("Similar Object already exist in storage")
	KeyNotFound         = StoreErr("Object with such key was not found in store")
	ValueSizeIsExceeded = StoreErr("Value size is more than MaxValueSize")

	StorageIsLocked = StoreErr("Could not lock store")
	ReceivedStop    = StoreErr("Received signal from ctx.Done()")
)

func (err StoreErr) Error() string {
	return string(err)
}
