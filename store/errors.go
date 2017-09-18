package store

type Err string

// To check is Err implements error interface
var _ = error(Err(""))

const (
	KeyAlreadyExist     = Err("Similar Object already exist in storage")
	KeyNotFound         = Err("Object with such key was not found in store")
	ValueSizeIsExceeded = Err("Value size is more than MaxValueSize")

	ReceivedStop = Err("Received signal from ctx.Done()")
)

func (err Err) Error() string {
	return string(err)
}
