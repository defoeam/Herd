package keyvaluestore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Snapshot struct {
	Data      map[string]json.RawMessage `json:"data"`
	Timestamp time.Time                  `json:"timestamp"`
}

func (kv *KeyValueStore) TakeSnapshot() error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	// Convert map[string][]byte to map[string]json.RawMessage
	convertedData := make(map[string]json.RawMessage)
	for k, v := range kv.data {
		convertedData[k] = json.RawMessage(v)
	}

	snapshot := Snapshot{
		Data:      convertedData,
		Timestamp: time.Now(),
	}

	snapshotData, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	snapshotFileName := fmt.Sprintf("snapshot_%s.json", snapshot.Timestamp.Format("20060102150405"))
	snapshotFile := filepath.Join(filepath.Dir(kv.logger.filename), snapshotFileName)

	if err := os.WriteFile(snapshotFile, snapshotData, 0600); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	// Truncate the transaction log
	if err := kv.logger.ClearLogs(); err != nil {
		return fmt.Errorf("failed to clear transaction logs: %w", err)
	}

	return nil
}

func (kv *KeyValueStore) LoadLatestSnapshot() error {
	snapshotDir := filepath.Dir(kv.logger.filename)
	snapshots, err := filepath.Glob(filepath.Join(snapshotDir, "snapshot_*.json"))
	if err != nil {
		return fmt.Errorf("failed to list snapshot files: %w", err)
	}

	if len(snapshots) == 0 {
		return nil // No snapshots found, which is fine
	}

	latestSnapshot := snapshots[len(snapshots)-1]
	snapshotData, err := os.ReadFile(latestSnapshot)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(snapshotData, &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data = make(map[string][]byte)
	for k, v := range snapshot.Data {
		kv.data[k] = []byte(v)
	}

	return nil
}
