package recording

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Header struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Env       map[string]string `json:"env,omitempty"`
}

type Recorder struct {
	mu        sync.Mutex
	file      *os.File
	startTime time.Time
	closed    bool
}

// NewRecorder creates a new asciinema cast v2 recorder
func NewRecorder(dir string, sessionID string, width, height int) (*Recorder, string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("%s.cast", sessionID))
	f, err := os.Create(path)
	if err != nil {
		return nil, "", err
	}

	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	header := Header{
		Version:   2,
		Width:     width,
		Height:    height,
		Timestamp: time.Now().Unix(),
		Env:       map[string]string{"TERM": "xterm-256color"},
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		f.Close()
		return nil, "", err
	}

	if _, err := f.Write(append(headerBytes, '\n')); err != nil {
		f.Close()
		return nil, "", err
	}

	return &Recorder{
		file:      f,
		startTime: time.Now(),
	}, path, nil
}

// WriteOutput records an output event
func (r *Recorder) WriteOutput(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed || r.file == nil {
		return len(p), nil
	}

	offset := time.Since(r.startTime).Seconds()
	event := []interface{}{offset, "o", string(p)}

	data, err := json.Marshal(event)
	if err != nil {
		return 0, err
	}

	_, err = r.file.Write(append(data, '\n'))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Resize records a terminal resize event
func (r *Recorder) Resize(width, height int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed || r.file == nil {
		return nil
	}

	offset := time.Since(r.startTime).Seconds()
	event := []interface{}{offset, "r", fmt.Sprintf("%dx%d", width, height)}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = r.file.Write(append(data, '\n'))
	return err
}

// Close closes the recording file
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

type OutputWriter struct {
	recorder *Recorder
}

// NewOutputWriter creates an io.Writer interface for the recorder
func NewOutputWriter(recorder *Recorder) io.Writer {
	return &OutputWriter{recorder: recorder}
}

func (w *OutputWriter) Write(p []byte) (n int, err error) {
	if w.recorder == nil {
		return len(p), nil
	}
	return w.recorder.WriteOutput(p)
}
