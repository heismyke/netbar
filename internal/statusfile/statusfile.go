package statusfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/heismyke/netbar/internal/monitor"
)

// DefaultPath returns the path used to remember the last network status.
func DefaultPath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "netbar", "status")
	}

	return filepath.Join(dir, "netbar", "status")
}

// Read reads a status value from path. Invalid or missing values return empty.
func Read(path string) monitor.Status {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	switch monitor.Status(strings.TrimSpace(string(data))) {
	case monitor.StatusOnline:
		return monitor.StatusOnline
	case monitor.StatusDegraded:
		return monitor.StatusDegraded
	case monitor.StatusOffline:
		return monitor.StatusOffline
	default:
		return ""
	}
}

// Write stores status at path. It is best-effort because status history should
// never break the foreground terminal session.
func Write(path string, status monitor.Status) {
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	_ = os.WriteFile(path, []byte(string(status)+"\n"), 0o644)
}
