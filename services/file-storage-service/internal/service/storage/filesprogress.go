package storage

import (
	"context"
	"sync"
)

type FileProgress struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type filesInProgress struct {
	files map[string]*FileProgress
	mu    sync.RWMutex
}

func newFilesInProgress() *filesInProgress {
	return &filesInProgress{
		files: make(map[string]*FileProgress),
	}
}

func (f *filesInProgress) get(key string) (*FileProgress, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	process, ok := f.files[key]
	if !ok {
		return nil, false
	}
	return process, true
}

func (f *filesInProgress) set(key string, value *FileProgress) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.files[key] = value
}

func (f *filesInProgress) delete(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.files, key)
}
