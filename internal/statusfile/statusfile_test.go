package statusfile

import (
	"path/filepath"
	"testing"

	"github.com/heismyke/netbar/internal/monitor"
)

func TestReadWriteStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "netbar", "status")

	Write(path, monitor.StatusOffline)

	status := Read(path)
	if status != monitor.StatusOffline {
		t.Fatalf("expected %q, got %q", monitor.StatusOffline, status)
	}
}

func TestReadMissingStatus(t *testing.T) {
	status := Read(filepath.Join(t.TempDir(), "missing"))
	if status != "" {
		t.Fatalf("expected empty status, got %q", status)
	}
}
