package main

import (
	"fmt"
	"github.io/uberate/sraft/cmd"
	"github.io/uberate/sraft/pkg/sraft"
	"os"
)

func main() {
	sc := sraft.DefaultServerConfig()

	if err := cmd.ReadConfig(
		"", "", "",
		true, "",
		sc); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s, err := sc.Server()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = s.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
