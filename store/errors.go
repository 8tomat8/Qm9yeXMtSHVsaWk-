package store

type StoreErr string

// To check is StoreErr implements error interface
var _ = error(StoreErr(""))

const (
	KeyAlreadyExist     = StoreErr("Similar Object already exist in storage")
	KeyNotFound         = StoreErr("Object with such key was not found in store")
	ValueSizeIsExceeded = StoreErr("Value size is more than MaxValueSize")

	ReceivedStop = StoreErr("Received signal from ctx.Done()")
)

func (err StoreErr) Error() string {
	return string(err)
}
