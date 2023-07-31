package secret

type loader interface {
	load() ([]byte, error)
}
