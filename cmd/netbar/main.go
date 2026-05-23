package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/heismyke/netbar/internal/monitor"
)

const (
	defaultHost     = "8.8.8.8:53"
	defaultInterval = 3 * time.Second
)

func main() {
	host := flag.String("host", defaultHost, "TCP host to probe")
	interval := flag.Duration("interval", defaultInterval, "connectivity check interval")
	once := flag.Bool("once", false, "run one check and exit")
	format := flag.String("format", "plain", "output format: plain or tmux")
	stateFile := flag.String("state-file", defaultStateFile(), "path used to remember the previous status")
	flag.Parse()

	cm := monitor.NewConnectivityManager(*host, *interval)
	previous := readStatus(*stateFile)

	if *once {
		result := cm.Check()
		printResult(result, previous, *format)
		writeStatus(*stateFile, result.Status)
		return
	}

	if err := cm.Start(); err != nil {
		log.Fatalf("start connectivity monitor: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	for {
		select {
		case result, ok := <-cm.Results():
			if !ok {
				return
			}

			printResult(result, previous, *format)
			writeStatus(*stateFile, result.Status)
			previous = result.Status
		case <-signals:
			cm.Stop()
		}
	}
}

func printResult(result monitor.Result, previous monitor.Status, format string) {
	switch format {
	case "plain":
		fmt.Println(formatPlain(result, previous))
	case "tmux":
		fmt.Println(formatTmux(result, previous))
	default:
		log.Fatalf("unsupported format %q", format)
	}
}

func formatPlain(result monitor.Result, previous monitor.Status) string {
	label := statusLabel(result.Status, previous)
	checkedAt := result.CheckedAt.Format(time.RFC3339)

	if result.Err != nil {
		return fmt.Sprintf("%s host=%s checked_at=%s error=%v", label, result.Host, checkedAt, result.Err)
	}

	return fmt.Sprintf("%s host=%s latency=%s checked_at=%s", label, result.Host, result.Latency.Round(time.Millisecond), checkedAt)
}

func formatTmux(result monitor.Result, previous monitor.Status) string {
	label := statusLabel(result.Status, previous)

	switch result.Status {
	case monitor.StatusOnline:
		return fmt.Sprintf("#[fg=white,bg=colour34,bold] %s #[default]", label)
	case monitor.StatusDegraded:
		return fmt.Sprintf("#[fg=black,bg=colour178,bold] %s #[default]", label)
	default:
		return fmt.Sprintf("#[fg=white,bg=colour160,bold] %s #[default]", label)
	}
}

func statusLabel(status monitor.Status, previous monitor.Status) string {
	if status == monitor.StatusOnline && (previous == monitor.StatusOffline || previous == monitor.StatusDegraded) {
		return "Back online"
	}

	switch status {
	case monitor.StatusOnline:
		return "Online"
	case monitor.StatusDegraded:
		return "Degraded"
	case monitor.StatusOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

func readStatus(path string) monitor.Status {
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

func writeStatus(path string, status monitor.Status) {
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	_ = os.WriteFile(path, []byte(string(status)+"\n"), 0o644)
}

func defaultStateFile() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "netbar", "status")
	}

	return filepath.Join(dir, "netbar", "status")
}
