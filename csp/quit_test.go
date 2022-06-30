package main

import (
	"testing"
	"time"
)

func TestConcurrentShutdown(t *testing.T) {
	f1 := shutdownMaker(1)
	f2 := shutdownMaker(8)

	err := ConcurrentShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err != nil {
		t.Errorf("want nil, actual: %s", err)
		return
	}

	err = ConcurrentShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err == nil {
		t.Error("want timeout, actual nil")
		return
	}

}

func TestSequentialShutdown(t *testing.T) {
	f1 := shutdownMaker(1)
	f2 := shutdownMaker(8)

	err := SequentialShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err != nil {
		t.Errorf("want nil, actual: %s", err)
		return
	}

	err = SequentialShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err == nil {
		t.Error("want timeout, actual nil")
		return
	}
}
