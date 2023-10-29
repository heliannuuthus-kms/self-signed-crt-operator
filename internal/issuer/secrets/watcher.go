package secrets

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/heliannuuthus/privateca-issuer/internal/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	"os"
	"path/filepath"
)

type secretWatcher interface {
	secretLoader
	Watch(string) error
	Notify()
}

type FileSecretWatcher struct {
	secretLoader
	watcher *fsnotify.Watcher
	files   cmap.ConcurrentMap[string, []string]
}

func NewFileSecretWatcher(loader secretLoader) (*FileSecretWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	self := &FileSecretWatcher{secretLoader: loader, watcher: watcher, files: cmap.New[[]string]()}

	go self.Notify()

	return self, nil

}

func (fsw *FileSecretWatcher) Watch(files ...string) error {
	var (
		done []string
		err  error
		fi   os.FileInfo
	)

	for _, p := range files {
		fi, err = os.Lstat(p)
		if err != nil {
			break
		}
		if fi.IsDir() {
			return fmt.Errorf("%q is a directory, not a file", p)
		}

		dir := filepath.Dir(p)
		if subFiles, ok := fsw.files.Get(dir); ok {
			subFiles = append(subFiles, p)
			fsw.files.Set(dir, subFiles)
			continue
		}
		err = fsw.watcher.Add(dir)
		if err != nil {
			break
		}
		done = append(done, p)

	}

	if err != nil {
		for _, f := range done {
			dir := filepath.Dir(f)
			if subFiles, ok := fsw.files.Get(dir); ok && len(subFiles) == 0 {
				err := fsw.watcher.Remove(dir)
				if err != nil {
					utils.Logger.V(4).Error(err, "remove file watch signal failed, path: %s", f)
				}
			}
		}
		return err
	}
	return nil
}

func (fsw *FileSecretWatcher) Notify() {
	for {
		select {
		// Read from Errors.
		case err, ok := <-fsw.watcher.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			utils.Logger.V(4).Error(err, "remove file watch signal failed, path: %s", f)
		// Read from Events.
		case e, ok := <-fsw.watcher.Events:
			if !ok {
				return
			}
			if e.Op&fsnotify.Create == fsnotify.Create || e.Op&fsnotify.Write == fsnotify.Write {
				file, err := os.ReadFile(e.Name)
				if err != nil {
					utils.Logger.V(4).Error(err, "read file failed", "file", e.Name, "state", e.Op)
					continue
				}
				err = fsw.Update(e.Name, file)
				if err != nil {
					utils.Logger.V(4).Error(err, "watching file update failed", "file", e.Name, "state", e.Op)
					continue
				}
			} else if e.Op&fsnotify.Rename == fsnotify.Rename {
				file, err := os.ReadFile(e.Name)
				if err != nil {
					utils.Logger.V(4).Error(err, "read file failed failed", "file", e.Name, "state", e.Op)
					continue
				}
				err = fsw.Update(e.Name, file)
				if err != nil {
					utils.Logger.V(4).Error(err, "watching file load failed", "file", e.Name, "state", e.Op)
				}
				err = fsw.watcher.Add(filepath.Dir(e.Name))
				if err != nil {
					utils.Logger.V(4).Error(err, "file rename re-watching file failed", "file_name", e.Name, "state", e.Op)
					continue
				}
			} else if e.Op&fsnotify.Remove == fsnotify.Remove {
				utils.Logger.V(0).Info("file(%s) remove down watch", e.Name)
			}
		}
	}
}
