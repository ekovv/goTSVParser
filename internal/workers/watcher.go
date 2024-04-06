package workers

import (
	"context"
	"fmt"
	"goTSVParser/config"
	"goTSVParser/internal/shema"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func NewWatcher(c config.Config) *Watcher {
	return &Watcher{timer: c.RefreshInterval, fromDir: c.DirectoryFrom, files: make(map[string]struct{})}
}

type Watcher struct {
	mutex   sync.RWMutex
	timer   int
	fromDir string
	files   map[string]struct{}
}

func (w *Watcher) InitCheckedFiles(files []shema.ParsedFiles) {
	for _, file := range files {
		w.files[file.File] = struct{}{}
	}
}

// Scan main scan directory
func (s *Watcher) Scan(ctx context.Context, out chan string) {
	go func() {
		timer := time.NewTicker(time.Duration(s.timer) * time.Second)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				err := filepath.Walk(s.fromDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						s.mutex.Lock()
						_, ok := s.files[path]
						if ok {
							s.mutex.Unlock()
							return nil
						} else {
							s.files[path] = struct{}{}
							s.mutex.Unlock()
						}
						out <- path
					}
					return nil
				})
				if err != nil {
					fmt.Errorf("error reading %w", err)
					return
				}
			}
		}
	}()
}
