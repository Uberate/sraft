package sraft

import (
	"fmt"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"strings"
)

// Command save the Value change info.
type Command struct {
	FullPath string
	Value    string
	Action   string
}

func (cmd *Command) Exec(storage storage.Storage) error {
	action := strings.ToLower(cmd.Action)
	switch action {
	case "put":
		storage.Put(cmd.FullPath, cmd.Value)
	case "delete":
		storage.Delete(cmd.FullPath)
	case "set":
		storage.Set(cmd.FullPath, cmd.Value)
	default:
		return fmt.Errorf("Unknown Action: %s, support put, delete, set. ", cmd.Action)
	}
	return nil
}
