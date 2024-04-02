package workers

import (
	"fmt"
	"goTSVParser/config"
	"goTSVParser/internal/shema"
	"os"
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

func (s *Watcher) Scan(out chan string) {
	go func() {
		timer := time.NewTicker(time.Duration(s.timer) * time.Second)

		defer timer.Stop()

		for range timer.C {
			filesFromDir, err := os.ReadDir(s.fromDir)
			if err != nil {
				fmt.Errorf("error reading %w", err)
				return
			}
			for _, file := range filesFromDir {
				if !file.IsDir() {
					s.mutex.Lock()
					_, ok := s.files[file.Name()]
					if ok {
						s.mutex.Unlock()
						continue
					} else {
						s.files[file.Name()] = struct{}{}
						s.mutex.Unlock()
					}
					out <- file.Name()

				} else {

				}
			}
		}
	}()

}
