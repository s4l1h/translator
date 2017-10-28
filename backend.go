package translator


// Backend backend interface
type Backend interface {
	Load() []*Translation
	Save(*Translation)
	Delete(*Translation)
}