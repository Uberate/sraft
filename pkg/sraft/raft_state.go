package sraft

import (
	"sync"
)

type StateMachine struct {
	logLock *sync.RWMutex

	Logs          []Log
	NextAppends   map[string]uint64
	NextCommitted uint64
}

func (sm *StateMachine) Init() {
	sm.Logs = make([]Log, 0)
	sm.NextAppends = map[string]uint64{}
	sm.NextCommitted = 0
}

func (sm *StateMachine) AppendOne(log Log) {
	sm.logLock.Lock()
	defer sm.logLock.Unlock()

	sm.Logs = append(sm.Logs, log)
}

func (sm *StateMachine) Commit(index uint64) {
	sm.logLock.Lock()
	defer sm.logLock.Unlock()
	sm.NextCommitted = index
}

func (sm *StateMachine) UpdateNextAppend(id string, nextAppend uint64) {
	sm.logLock.Lock()
	defer sm.logLock.Unlock()

	sm.NextAppends[id] = nextAppend
}

func (sm *StateMachine) GetNextAppendIndex(id string) uint64 {
	sm.logLock.RLock()
	defer sm.logLock.RUnlock()

	if value, ok := sm.NextAppends[id]; ok {
		return value
	} else {
		return 0
	}
}
