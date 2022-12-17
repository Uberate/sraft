package storage

import (
	"github.io/uberate/sraft/pkg/sraft"
	"regexp"
	"sync"
)

const MemoryV1EngineName = "memory_v1"

type MemoryV1Engine struct {
	innerMap sync.Map
}

func (m *MemoryV1Engine) Name() string {
	return MemoryV1EngineName
}

func (m *MemoryV1Engine) SetConfig(config sraft.AnyConfig) error {
	m.innerMap = sync.Map{}

	return nil
}

func (m *MemoryV1Engine) Size() uint64 {
	res := uint64(0)
	m.innerMap.Range(func(key, value any) bool {
		kStr := key.(string)
		vStr := value.(string)

		res += uint64(len([]byte(kStr)))
		res += uint64(len([]byte(vStr)))

		return true
	})

	return res
}

func (m *MemoryV1Engine) Len() uint64 {
	count := uint64(0)

	m.innerMap.Range(func(key, value any) bool {
		count++
		return true
	})

	return count
}

func (m *MemoryV1Engine) Paths(pattern string) []string {
	var res []string

	if len(pattern) == 0 {
		m.innerMap.Range(func(key, value any) bool {
			kStr := key.(string)
			res = append(res, kStr)
			return true
		})
	} else {
		rp, err := regexp.Compile(pattern)
		if err != nil {
			return []string{}
		}

		m.innerMap.Range(func(key, value any) bool {
			kStr := key.(string)
			if rp.MatchString(kStr) {
				res = append(res, kStr)
			}
			return true
		})
	}

	return res
}

func (m *MemoryV1Engine) Get(path string) (string, bool) {
	value, ok := m.innerMap.Load(path)
	if !ok {
		return "", false
	}
	return value.(string), ok
}

func (m *MemoryV1Engine) Put(path, value string) bool {
	if !m.ContainPath(path) {
		m.Set(path, value)
		return true
	} else {
		return false
	}
}

func (m *MemoryV1Engine) Set(path, value string) {
	m.innerMap.Store(path, value)
}

func (m *MemoryV1Engine) Delete(path string) (string, bool) {
	value, ok := m.innerMap.LoadAndDelete(path)
	if !ok {
		return "", false
	}
	return value.(string), ok
}

func (m *MemoryV1Engine) ContainPath(path string) bool {
	_, ok := m.innerMap.Load(path)
	return ok
}

func (m *MemoryV1Engine) Clean() error {
	m.innerMap = sync.Map{}

	return nil
}
