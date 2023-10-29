package secrets

import (
	"github.com/fsnotify/fsnotify"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type secretLoader interface {
	Load(string) ([]byte, error)
	Update(string, []byte) error
	Offload(string) error
}

type FileSecretLoader struct {
	cache cmap.ConcurrentMap[string, []byte]
}

func NewFileSecretLoader() *FileSecretLoader {
	return &FileSecretLoader{cache: cmap.New[[]byte]()}
}

func (fsl *FileSecretLoader) Load(path string) ([]byte, error) {

	if keys, ok := fsl.cache.Get(path); ok {
		return keys, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {

		return nil, err
	}
}

func (fsl *FileSecretLoader) Offload(path string) ([]byte, error) {

	if keys, ok := fsl.cache.Get(path); ok {
		return keys, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {

		return nil, err
	}
}
