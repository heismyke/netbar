package monitor

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestRunChecksReturnsOnline(t *testing.T) {
	manager := testManager(t)
	manager.degradedThreshold = time.Hour

	result := manager.runChecks()

	if result.Status != StatusOnline {
		t.Fatalf("expected status %q, got %q", StatusOnline, result.Status)
	}
	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}
}

func TestRunChecksReturnsDegraded(t *testing.T) {
	manager := testManager(t)
	manager.degradedThreshold = -time.Nanosecond

	result := manager.runChecks()

	if result.Status != StatusDegraded {
		t.Fatalf("expected status %q, got %q", StatusDegraded, result.Status)
	}
}

func TestRunChecksReturnsOfflineWhenDNSFails(t *testing.T) {
	manager := testManager(t)
	manager.lookupHost = func(context.Context, string) ([]string, error) {
		return nil, errors.New("lookup failed")
	}

	result := manager.runChecks()

	if result.Status != StatusOffline {
		t.Fatalf("expected status %q, got %q", StatusOffline, result.Status)
	}
	if result.Err == nil {
		t.Fatal("expected error")
	}
}

func TestStartStopClosesResults(t *testing.T) {
	manager := testManager(t)
	manager.interval = time.Hour

	if err := manager.Start(); err != nil {
		t.Fatalf("start manager: %v", err)
	}

	results := manager.Results()
	select {
	case <-results:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for initial result")
	}

	manager.Stop()
	manager.Stop()

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("expected results channel to close")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for results channel to close")
	}
}

func TestStartRejectsInvalidConfiguration(t *testing.T) {
	manager := testManager(t)
	manager.host = ""

	if err := manager.Start(); err == nil {
		t.Fatal("expected empty host to be rejected")
	}

	manager = testManager(t)
	manager.interval = 0

	if err := manager.Start(); err == nil {
		t.Fatal("expected zero interval to be rejected")
	}
}

func testManager(t *testing.T) *Manager {
	t.Helper()

	manager := NewConnectivityManager("127.0.0.1:53", time.Millisecond)
	manager.lookupHost = func(context.Context, string) ([]string, error) {
		return []string{"127.0.0.1"}, nil
	}
	manager.dialContext = func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		t.Cleanup(func() {
			client.Close()
			server.Close()
		})
		return client, nil
	}

	return manager
}
