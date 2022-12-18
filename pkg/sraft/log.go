package sraft

// Log is a part of sraft, not the log of system, log storage the Value change of status-machine.
type Log struct {
	Cmd      Command
	CreateBy string
	CreateAt string
	Term     string
}

type Logs struct {
	LogEntries []Logs

	LastAppendAt  uint64
	LastCommitted uint64
}
