package sraft

import "testing"

func TestLogger(t *testing.T) {
	_ = DefaultLogConfig().ToLogger()
	conf := DefaultLogConfig()
	conf.Level = "not allow level"
	_ = conf.ToLogger()
}
