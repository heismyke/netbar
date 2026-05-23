package monitor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	defaultProbeTimeout       = 2 * time.Second
	defaultDegradedThreshold  = 300 * time.Millisecond
	defaultDNSProbeHostname   = "dns.google"
	defaultResultBufferLength = 1
)

// Status describes the latest connectivity state.
type Status string

const (
	StatusOnline   Status = "online"
	StatusDegraded Status = "degraded"
	StatusOffline  Status = "offline"
)

// Result is a single connectivity check result.
type Result struct {
	Status    Status
	Host      string
	CheckedAt time.Time
	Latency   time.Duration
	Err       error
}

// Manager periodically checks network connectivity and publishes the latest
// result to subscribers.
type Manager struct {
	host              string
	interval          time.Duration
	degradedThreshold time.Duration
	probeTimeout      time.Duration
	dnsProbeHostname  string

	dialContext func(context.Context, string, string) (net.Conn, error)
	lookupHost  func(context.Context, string) ([]string, error)

	mu       sync.Mutex
	running  bool
	stopCh   chan struct{}
	resultCh chan Result
}

// NewConnectivityManager creates a connectivity manager for host.
func NewConnectivityManager(host string, interval time.Duration) *Manager {
	return &Manager{
		host:              host,
		interval:          interval,
		degradedThreshold: defaultDegradedThreshold,
		probeTimeout:      defaultProbeTimeout,
		dnsProbeHostname:  defaultDNSProbeHostname,
		dialContext:       (&net.Dialer{}).DialContext,
		lookupHost:        net.DefaultResolver.LookupHost,
		stopCh:            make(chan struct{}),
		resultCh:          make(chan Result, defaultResultBufferLength),
	}
}

// Results returns the stream of connectivity results.
func (m *Manager) Results() <-chan Result {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.resultCh
}

// Check runs one connectivity check immediately.
func (m *Manager) Check() Result {
	return m.runChecks()
}

// Start begins monitoring. Calling Start more than once is a no-op.
func (m *Manager) Start() error {
	if err := m.validate(); err != nil {
		return err
	}

	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return nil
	}

	m.running = true
	m.stopCh = make(chan struct{})
	m.resultCh = make(chan Result, defaultResultBufferLength)
	stopCh := m.stopCh
	resultCh := m.resultCh
	m.mu.Unlock()

	go m.run(stopCh, resultCh)
	return nil
}

func (m *Manager) validate() error {
	if m.host == "" {
		return errors.New("monitor host is required")
	}

	if m.interval <= 0 {
		return errors.New("monitor interval must be greater than zero")
	}

	return nil
}

// Stop stops monitoring. Calling Stop before Start or multiple times is safe.
func (m *Manager) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}

	close(m.stopCh)
	m.running = false
	m.mu.Unlock()
}

func (m *Manager) run(stopCh <-chan struct{}, resultCh chan Result) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	defer close(resultCh)

	publishLatest(resultCh, m.runChecks())

	for {
		select {
		case <-ticker.C:
			publishLatest(resultCh, m.runChecks())
		case <-stopCh:
			return
		}
	}
}

func publishLatest(resultCh chan Result, result Result) {
	select {
	case resultCh <- result:
		return
	default:
	}

	select {
	case <-resultCh:
	default:
	}

	select {
	case resultCh <- result:
	default:
	}
}

func (m *Manager) runChecks() Result {
	now := time.Now()

	if err := m.checkDNS(); err != nil {
		return Result{
			Status:    StatusOffline,
			Host:      m.host,
			CheckedAt: now,
			Err:       fmt.Errorf("dns check failed: %w", err),
		}
	}

	latency, err := m.checkTCP()
	if err != nil {
		return Result{
			Status:    StatusOffline,
			Host:      m.host,
			CheckedAt: now,
			Err:       fmt.Errorf("tcp check failed: %w", err),
		}
	}

	status := StatusOnline
	if latency > m.degradedThreshold {
		status = StatusDegraded
	}

	return Result{
		Status:    status,
		Host:      m.host,
		CheckedAt: now,
		Latency:   latency,
	}
}

func (m *Manager) checkDNS() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.probeTimeout)
	defer cancel()

	_, err := m.lookupHost(ctx, m.dnsProbeHostname)
	return err
}

func (m *Manager) checkTCP() (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.probeTimeout)
	defer cancel()

	start := time.Now()
	conn, err := m.dialContext(ctx, "tcp", m.host)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return time.Since(start), nil
}
