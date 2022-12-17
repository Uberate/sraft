package point

import (
	"github.io/uberate/sraft/pkg/sraft"
	"github.io/uberate/sraft/tests"
	"testing"
	"time"
)

func QuickHttpV1Server() *HttpV1Server {
	return &HttpV1Server{
		stopChan:    make(chan bool, 0),
		alreadyStop: make(chan bool, 0),
		logger:      tests.QuickTestLogger(),
		handlers:    map[string]sraft.Handler{},

		Point: "127.0.0.1:9867",
	}
}

func TestStop(t *testing.T) {
	runTime := 5 * time.Second

	serverRunBeginTime := time.Now()

	server := QuickHttpV1Server()

	go func() {
		err := server.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	for time.Now().Before(serverRunBeginTime.Add(runTime)) {
	}

	awaitTime := 3 * time.Second
	stopBeginTime := time.Now()
	closed := false
	go func() {
		if err := server.Stop(); err != nil {
			t.Error(err)
		}
		closed = true
	}()

	for !closed && time.Now().Before(stopBeginTime.Add(awaitTime)) {
	}

	if !closed {
		t.Errorf("Can't stop server!")
	}

	return
}
