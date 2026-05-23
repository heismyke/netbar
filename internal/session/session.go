package session

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/heismyke/netbar/internal/monitor"
	"github.com/heismyke/netbar/internal/statusfile"
	"golang.org/x/term"
)

// Config controls an interactive netbar session.
type Config struct {
	Host      string
	Interval  time.Duration
	StateFile string
	Command   []string
}

// Run starts an interactive command with a network status bar pinned to the
// bottom terminal row.
func Run(ctx context.Context, cfg Config) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return errors.New("interactive session requires a terminal")
	}

	command := cfg.Command
	if len(command) == 0 {
		command = defaultCommand()
	}

	rows, cols, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("get terminal size: %w", err)
	}
	if rows < 3 {
		return errors.New("terminal must be at least 3 rows tall")
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(), "NETBAR=1")

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(rows - 1),
		Cols: uint16(cols),
	})
	if err != nil {
		return fmt.Errorf("start command: %w", err)
	}
	defer ptmx.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("set terminal raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	renderer := &renderer{
		out:  os.Stdout,
		rows: rows,
		cols: cols,
	}
	renderer.enter()
	defer renderer.leave()

	manager := monitor.NewConnectivityManager(cfg.Host, cfg.Interval)
	if err := manager.Start(); err != nil {
		return err
	}
	defer manager.Stop()

	previous := statusfile.Read(cfg.StateFile)
	renderer.draw(monitor.Result{Status: previous, Host: cfg.Host})

	errCh := make(chan error, 2)
	waitCh := make(chan error, 1)

	go func() {
		_, err := io.Copy(ptmx, os.Stdin)
		errCh <- err
	}()

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				renderer.writeCommand(buf[:n])
			}
			if err != nil {
				if errors.Is(err, os.ErrClosed) || errors.Is(err, io.EOF) || errors.Is(err, syscall.EIO) {
					errCh <- nil
					return
				}
				errCh <- err
				return
			}
		}
	}()

	go func() {
		for result := range manager.Results() {
			renderer.drawTransition(result, previous)
			statusfile.Write(cfg.StateFile, result.Status)
			previous = result.Status
		}
	}()

	go func() {
		waitCh <- cmd.Wait()
	}()

	resizeSignals := make(chan os.Signal, 1)
	signal.Notify(resizeSignals, syscall.SIGWINCH)
	defer signal.Stop(resizeSignals)

	for {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Signal(syscall.SIGTERM)
			return ctx.Err()
		case err := <-waitCh:
			return err
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-resizeSignals:
			if err := renderer.resize(ptmx); err != nil {
				return err
			}
		}
	}
}

type renderer struct {
	mu      sync.Mutex
	out     *os.File
	rows    int
	cols    int
	current string
}

func (r *renderer) enter() {
	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Fprintf(r.out, "\x1b[?1049h\x1b[2J\x1b[1;%dr", r.rows-1)
}

func (r *renderer) leave() {
	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Fprintf(r.out, "\x1b[r\x1b[%d;1H\x1b[2K\x1b[?1049l", r.rows)
}

func (r *renderer) writeCommand(data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, _ = r.out.Write(data)
	if r.current != "" {
		r.drawLocked(r.current)
	}
}

func (r *renderer) draw(result monitor.Result) {
	r.drawTransition(result, "")
}

func (r *renderer) drawTransition(result monitor.Result, previous monitor.Status) {
	label := statusLabel(result.Status, previous)
	text := label
	if result.Status != "" {
		text = fmt.Sprintf("netbar: %s", label)
	}
	if result.Latency > 0 && result.Status != monitor.StatusOffline {
		text = fmt.Sprintf("%s  %s", text, result.Latency.Round(time.Millisecond))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.current = renderBar(result.Status, text, r.cols)
	r.drawLocked(r.current)
}

func (r *renderer) drawLocked(bar string) {
	fmt.Fprintf(r.out, "\x1b[s\x1b[%d;1H%s\x1b[u", r.rows, bar)
}

func (r *renderer) resize(ptmx *os.File) error {
	rows, cols, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}
	if rows < 3 {
		return errors.New("terminal must be at least 3 rows tall")
	}

	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(rows - 1),
		Cols: uint16(cols),
	}); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.rows = rows
	r.cols = cols
	fmt.Fprintf(r.out, "\x1b[1;%dr", r.rows-1)
	if r.current != "" {
		r.drawLocked(r.current)
	}

	return nil
}

func renderBar(status monitor.Status, text string, width int) string {
	if width <= 0 {
		width = 80
	}

	if text == "" {
		text = "netbar"
	}

	text = " " + text + " "
	if len(text) > width {
		text = text[:width]
	}

	left := (width - len(text)) / 2
	right := width - len(text) - left

	return statusColor(status) + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "\x1b[0m"
}

func statusColor(status monitor.Status) string {
	switch status {
	case monitor.StatusOnline:
		return "\x1b[1;37;42m"
	case monitor.StatusDegraded:
		return "\x1b[1;30;43m"
	case monitor.StatusOffline:
		return "\x1b[1;37;41m"
	default:
		return "\x1b[1;37;44m"
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
		return "Checking..."
	}
}

func defaultCommand() []string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	return []string{shell}
}
