package sraft

// Log is a part of sraft, not the log of system, log storage the Value change of status-machine.
type Log struct {
	Type     string
	Action   string
	Path     string
	Value    string
	Index    uint64
	CreateBy string
	CreateAt uint64
	Term     uint64
}

type Logs struct {
	LogEntries []Log

	NextAppend    uint64
	NextCommitted uint64
}

func (l *Logs) GetLastAppendTerm() uint64 {
	if l.NextAppend == 0 {
		return 0
	}
	return l.LogEntries[l.NextAppend-1].Term
}
