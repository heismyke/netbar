package session

import (
	"strings"
	"testing"

	"github.com/heismyke/netbar/internal/monitor"
)

func TestRenderBarFillsWidth(t *testing.T) {
	bar := renderBar(monitor.StatusOnline, "netbar: Online", 40)

	if !strings.Contains(bar, "netbar: Online") {
		t.Fatalf("expected bar to contain label, got %q", bar)
	}
	if !strings.Contains(bar, "\x1b[1;37;42m") {
		t.Fatalf("expected online color, got %q", bar)
	}
}

func TestRenderBarTruncatesLongText(t *testing.T) {
	bar := renderBar(monitor.StatusOffline, "this text is too long", 10)

	visible := strings.TrimSuffix(strings.TrimPrefix(bar, statusColor(monitor.StatusOffline)), "\x1b[0m")
	if len(visible) != 10 {
		t.Fatalf("expected visible bar length 10, got %d from %q", len(visible), visible)
	}
}

func TestSessionStatusLabelBackOnline(t *testing.T) {
	label := statusLabel(monitor.StatusOnline, monitor.StatusOffline)
	if label != "Back online" {
		t.Fatalf("expected Back online, got %q", label)
	}
}
