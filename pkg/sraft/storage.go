package sraft

// Storage is the abstract of the storage layout. And the implements not care the thread-safe. That is the high layout
// task.
type Storage interface {

	// Name of the storage implement.
	Name() string

	// SetConfig will set the config info of Storage implements.
	SetConfig(config AnyConfig) error

	// Size return the bytes of storage, if the Len() is zero, the Size should return zero.
	Size() uint64

	// Len return the elements of storage
	Len() uint64

	// Paths return the paths for pattern, if nil, return all path.
	Paths(pattern string) []string

	// Get return the value of specify path, if specify path not exists, return "" with false. Else return value with
	// true.
	Get(path string) (string, bool)

	// Put set the value when path not exists. If specify path exists, return false, else return true.
	Put(path string) bool

	// Set the value ignore path exists. It will cover the value of specify path.
	Set(path string)

	// Delete specify path, and if specify path not exists, return "" with false, else return value with true.
	Delete(path string) (string, bool)

	// ContainPath return true when specify path exists, else return false.
	ContainPath(path string) bool

	// Clean the values of storage.
	Clean() error
}
