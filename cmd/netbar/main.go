package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heismyke/netbar/internal/monitor"
	"github.com/heismyke/netbar/internal/session"
	"github.com/heismyke/netbar/internal/statusfile"
	"golang.org/x/term"
)

const (
	defaultHost     = "8.8.8.8:53"
	defaultInterval = 3 * time.Second
)

func main() {
	host := flag.String("host", defaultHost, "TCP host to probe")
	interval := flag.Duration("interval", defaultInterval, "connectivity check interval")
	once := flag.Bool("once", false, "run one check and exit")
	stream := flag.Bool("stream", false, "print continuous status updates instead of opening an interactive session")
	format := flag.String("format", "plain", "output format: plain or tmux")
	stateFile := flag.String("state-file", statusfile.DefaultPath(), "path used to remember the previous status")
	flag.Parse()

	cm := monitor.NewConnectivityManager(*host, *interval)
	previous := statusfile.Read(*stateFile)

	if *once {
		result := cm.Check()
		printResult(result, previous, *format)
		statusfile.Write(*stateFile, result.Status)
		return
	}

	if !*stream && term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := session.Run(ctx, session.Config{
			Host:      *host,
			Interval:  *interval,
			StateFile: *stateFile,
			Command:   flag.Args(),
		}); err != nil && err != context.Canceled {
			log.Fatalf("run session: %v", err)
		}
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
			statusfile.Write(*stateFile, result.Status)
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
