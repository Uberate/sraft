package storage

import (
	"github.io/uberate/sraft/pkg/plugins"
	"testing"
)

func TestMemoryV1Common(t *testing.T) {
	s, ok := GetStorageEngine(MemoryV1EngineName)
	if !ok {
		t.Errorf("Not register the %s", MemoryV1EngineName)
		return
	}
	_ = s.SetConfig(plugins.AnyConfig{})

	CommonStorageTest(s, t)
}
