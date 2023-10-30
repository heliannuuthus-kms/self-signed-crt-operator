package secrets

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"os"
)

type secretLoader interface {
	Load(string) ([]byte, error)
	Update(string, []byte) error
	Offload(string) error
}

type FileSecretLoader struct {
	secretWatcher
	cache cmap.ConcurrentMap[string, []byte]
}

func NewFileSecretLoader() (secretLoader, error) {
	loader := &FileSecretLoader{cache: cmap.New[[]byte]()}
	watcher, err := NewFileSecretWatcher(loader.Update)
	if err != nil {
		return nil, err
	}
	loader.secretWatcher = watcher
	return loader, nil
}

func (fsl *FileSecretLoader) Load(path string) ([]byte, error) {

	if keys, ok := fsl.cache.Get(path); ok {
		return keys, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = fsl.Watch(path)

	if err != nil {
		return nil, err
	}

	fsl.cache.Set(path, file)

	return file, nil
}

func (fsl *FileSecretLoader) Offload(path string) error {

	fsl.cache.Remove(path)
	return nil
}

func (fsl *FileSecretLoader) Update(path string, keys []byte) error {

	fsl.cache.Set(path, keys)

	return nil
}
