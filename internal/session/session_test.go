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
	if !strings.Contains(bar, "\x1b[1;37;48;5;28m") {
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

func TestChildEnvOverridesTerminalSize(t *testing.T) {
	env := childEnv([]string{
		"PATH=/bin",
		"NETBAR=old",
		"LINES=10",
		"COLUMNS=20",
	}, 40, 120)

	joined := strings.Join(env, "\n")
	for _, want := range []string{"PATH=/bin", "NETBAR=1", "LINES=40", "COLUMNS=120"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected env to include %q, got %q", want, joined)
		}
	}
	if strings.Contains(joined, "NETBAR=old") || strings.Contains(joined, "LINES=10") || strings.Contains(joined, "COLUMNS=20") {
		t.Fatalf("expected stale terminal env to be replaced, got %q", joined)
	}
}

func TestTerminalSizeUsesOverrides(t *testing.T) {
	rows, cols, err := terminalSize(40, 120)
	if err != nil {
		t.Fatalf("expected overrides to avoid terminal size error, got %v", err)
	}
	if rows != 40 || cols != 120 {
		t.Fatalf("expected 40x120, got %dx%d", rows, cols)
	}
}
