package storage

import (
	"testing"
)

func TestMemoryV1Common(t *testing.T) {
	s, ok := GetStorageEngine(MemoryV1EngineName)
	if !ok {
		t.Errorf("Not register the %s", MemoryV1EngineName)
		return
	}

	CommonStorageTest(s, t)
}
