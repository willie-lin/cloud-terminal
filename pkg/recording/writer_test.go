package recording

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func TestRecorder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "recording_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sessionID := "test_sess_123"
	rec, path, err := NewRecorder(tempDir, sessionID, 100, 30)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	testMessage := "Hello Terminal\r\n"
	_, err = rec.WriteOutput([]byte(testMessage))
	if err != nil {
		t.Fatalf("failed to write output: %v", err)
	}

	if err := rec.Close(); err != nil {
		t.Fatalf("failed to close recorder: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("recording file does not exist: %s", path)
	}

	// Read and verify header
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatalf("empty recording file")
	}

	var header Header
	if err := json.Unmarshal(scanner.Bytes(), &header); err != nil {
		t.Fatalf("failed to parse header: %v", err)
	}

	if header.Version != 2 || header.Width != 100 || header.Height != 30 {
		t.Errorf("unexpected header values: %+v", header)
	}

	// Read event line
	if !scanner.Scan() {
		t.Fatalf("missing event line")
	}

	var event []interface{}
	if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
		t.Fatalf("failed to parse event: %v", err)
	}

	if len(event) < 3 || event[1] != "o" || event[2] != testMessage {
		t.Errorf("unexpected event structure: %+v", event)
	}
}
