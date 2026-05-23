package main

import (
	"strings"
	"testing"
	"time"

	"github.com/heismyke/netbar/internal/monitor"
)

func TestStatusLabelBackOnline(t *testing.T) {
	label := statusLabel(monitor.StatusOnline, monitor.StatusOffline)
	if label != "Back online" {
		t.Fatalf("expected Back online, got %q", label)
	}
}

func TestFormatTmux(t *testing.T) {
	result := monitor.Result{
		Status:    monitor.StatusOffline,
		Host:      "8.8.8.8:53",
		CheckedAt: time.Now(),
	}

	output := formatTmux(result, "")
	if !strings.Contains(output, "Offline") {
		t.Fatalf("expected tmux output to include Offline, got %q", output)
	}
	if !strings.Contains(output, "#[") {
		t.Fatalf("expected tmux styling, got %q", output)
	}
}

func TestFormatPlain(t *testing.T) {
	result := monitor.Result{
		Status:    monitor.StatusOnline,
		Host:      "8.8.8.8:53",
		CheckedAt: time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC),
		Latency:   42 * time.Millisecond,
	}

	output := formatPlain(result, "")
	if !strings.Contains(output, "Online") {
		t.Fatalf("expected plain output to include Online, got %q", output)
	}
	if !strings.Contains(output, "latency=42ms") {
		t.Fatalf("expected plain output to include latency, got %q", output)
	}
}
