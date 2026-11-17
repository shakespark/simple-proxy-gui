package proxy

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultHistoryPoints = 180
	defaultMaxLogs       = 2000
)

// TrafficPoint represents one 1-second bandwidth sample for the traffic chart.
type TrafficPoint struct {
	Timestamp int64  `json:"timestamp"`
	Label     string `json:"label"`
	UpBps     int64  `json:"up_bps"`
	DownBps   int64  `json:"down_bps"`
}

// AccessRecord represents a single proxied request or tunnel session.
type AccessRecord struct {
	ID         uint64 `json:"id"`
	Time       string `json:"time"`
	ClientIP   string `json:"client_ip"`
	Method     string `json:"method"`
	Target     string `json:"target"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	UpBytes    int64  `json:"up_bytes"`
	DownBytes  int64  `json:"down_bytes"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

// MonitorSnapshot is a consistent point-in-time view returned to the UI.
type MonitorSnapshot struct {
	TotalUpBytes    int64          `json:"total_up_bytes"`
	TotalDownBytes  int64          `json:"total_down_bytes"`
	CurrentUpBps    int64          `json:"current_up_bps"`
	CurrentDownBps  int64          `json:"current_down_bps"`
	ActiveConns     int64          `json:"active_conns"`
	TotalRequests   int64          `json:"total_requests"`
	BlockedRequests int64          `json:"blocked_requests"`
	History         []TrafficPoint `json:"history"`
	Logs            []AccessRecord `json:"logs"`
}

// AccessEntry tracks live counters and completion metadata for an active or completed request.
type AccessEntry struct {
	id         uint64
	startedAt  time.Time
	timeStr    string
	clientIP   string
	method     string
	target     string
	upBytes    atomic.Int64
	downBytes  atomic.Int64
	mu         sync.RWMutex
	status     string
	statusCode int
	durationMs int64
	errMsg     string
	monitor    *Monitor
}

// AddUp increments uploaded bytes (client -> target) for both this entry and global totals.
func (e *AccessEntry) AddUp(n int64) {
	if e == nil || n <= 0 {
		return
	}
	e.upBytes.Add(n)
	if e.monitor != nil {
		e.monitor.totalUpBytes.Add(n)
	}
}

// AddDown increments downloaded bytes (target -> client) for both this entry and global totals.
func (e *AccessEntry) AddDown(n int64) {
	if e == nil || n <= 0 {
		return
	}
	e.downBytes.Add(n)
	if e.monitor != nil {
		e.monitor.totalDownBytes.Add(n)
	}
}

// Finish marks the request/tunnel as completed and records its final status and duration.
func (e *AccessEntry) Finish(status string, statusCode int, errMsg string) {
	if e == nil {
		return
	}
	dur := time.Since(e.startedAt).Milliseconds()
	e.mu.Lock()
	e.status = status
	e.statusCode = statusCode
	e.durationMs = dur
	e.errMsg = errMsg
	e.mu.Unlock()
}

func (e *AccessEntry) snapshot() AccessRecord {
	e.mu.RLock()
	status := e.status
	statusCode := e.statusCode
	durationMs := e.durationMs
	errMsg := e.errMsg
	e.mu.RUnlock()

	if status == "active" {
		durationMs = time.Since(e.startedAt).Milliseconds()
	}
	return AccessRecord{
		ID:         e.id,
		Time:       e.timeStr,
		ClientIP:   e.clientIP,
		Method:     e.method,
		Target:     e.target,
		Status:     status,
		StatusCode: statusCode,
		UpBytes:    e.upBytes.Load(),
		DownBytes:  e.downBytes.Load(),
		DurationMs: durationMs,
		Error:      errMsg,
	}
}

// Monitor collects real-time traffic throughput and recent access records.
type Monitor struct {
	totalUpBytes    atomic.Int64
	totalDownBytes  atomic.Int64
	activeConns     atomic.Int64
	totalRequests   atomic.Int64
	blockedRequests atomic.Int64
	nextID          atomic.Uint64

	maxLogs int

	mu             sync.RWMutex
	history        []TrafficPoint
	logs           []*AccessEntry
	lastSampleUp   int64
	lastSampleDown int64
	lastSampleAt   time.Time

	stopCh chan struct{}
	once   sync.Once
}

// NewMonitor initializes a Monitor and starts its 1-second traffic sampler.
func NewMonitor(maxLogs int) *Monitor {
	if maxLogs <= 0 {
		maxLogs = defaultMaxLogs
	}
	now := time.Now()
	m := &Monitor{
		maxLogs:      maxLogs,
		history:      make([]TrafficPoint, 0, defaultHistoryPoints),
		logs:         make([]*AccessEntry, 0, 256),
		lastSampleAt: now,
		stopCh:       make(chan struct{}),
	}
	// Pre-fill chart baseline so the graph renders smoothly from second zero.
	for i := defaultHistoryPoints - 1; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Second)
		m.history = append(m.history, TrafficPoint{
			Timestamp: t.Unix(),
			Label:     t.Format("15:04:05"),
		})
	}
	go m.runSampler()
	return m
}

// Close stops the background sampler goroutine.
func (m *Monitor) Close() {
	if m == nil {
		return
	}
	m.once.Do(func() {
		close(m.stopCh)
	})
}

func (m *Monitor) runSampler() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case now := <-ticker.C:
			m.SampleAt(now)
		}
	}
}

// SampleAt records a throughput point using deltas since the previous sample.
func (m *Monitor) SampleAt(now time.Time) {
	upNow := m.totalUpBytes.Load()
	downNow := m.totalDownBytes.Load()

	m.mu.Lock()
	defer m.mu.Unlock()

	elapsed := now.Sub(m.lastSampleAt).Seconds()
	if elapsed <= 0 {
		elapsed = 1
	}
	upDelta := max(int64(0), upNow-m.lastSampleUp)
	downDelta := max(int64(0), downNow-m.lastSampleDown)

	pt := TrafficPoint{
		Timestamp: now.Unix(),
		Label:     now.Format("15:04:05"),
		UpBps:     int64(float64(upDelta) / elapsed),
		DownBps:   int64(float64(downDelta) / elapsed),
	}

	m.lastSampleUp = upNow
	m.lastSampleDown = downNow
	m.lastSampleAt = now

	if len(m.history) >= defaultHistoryPoints {
		copy(m.history, m.history[1:])
		m.history[len(m.history)-1] = pt
	} else {
		m.history = append(m.history, pt)
	}
}

// RecordBlocked records an immediately rejected request (whitelist block or auth failure).
func (m *Monitor) RecordBlocked(clientIP, method, target, status string, statusCode int, errMsg string) {
	if m == nil {
		return
	}
	m.totalRequests.Add(1)
	if status == "blocked" {
		m.blockedRequests.Add(1)
	}
	now := time.Now()
	entry := &AccessEntry{
		id:         m.nextID.Add(1),
		startedAt:  now,
		timeStr:    now.Format("2006-01-02 15:04:05"),
		clientIP:   clientIP,
		method:     method,
		target:     target,
		status:     status,
		statusCode: statusCode,
		errMsg:     errMsg,
		monitor:    m,
	}
	m.appendLog(entry)
}

// StartRequest creates an active AccessEntry for a permitted request or tunnel.
func (m *Monitor) StartRequest(clientIP, method, target string) *AccessEntry {
	if m == nil {
		return nil
	}
	m.totalRequests.Add(1)
	now := time.Now()
	entry := &AccessEntry{
		id:         m.nextID.Add(1),
		startedAt:  now,
		timeStr:    now.Format("2006-01-02 15:04:05"),
		clientIP:   clientIP,
		method:     method,
		target:     target,
		status:     "active",
		statusCode: 0,
		monitor:    m,
	}
	m.appendLog(entry)
	return entry
}

// IncActiveConns adjusts the active connection gauge.
func (m *Monitor) IncActiveConns(delta int64) {
	if m == nil {
		return
	}
	m.activeConns.Add(delta)
}

func (m *Monitor) appendLog(entry *AccessEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.logs) >= m.maxLogs {
		drop := len(m.logs) - m.maxLogs + 1
		copy(m.logs, m.logs[drop:])
		m.logs = m.logs[:m.maxLogs-1]
	}
	m.logs = append(m.logs, entry)
}

// SetMaxLogs updates the maximum retained log entries.
func (m *Monitor) SetMaxLogs(maxLogs int) {
	if m == nil || maxLogs <= 0 {
		return
	}
	m.mu.Lock()
	m.maxLogs = maxLogs
	if len(m.logs) > maxLogs {
		drop := len(m.logs) - maxLogs
		copy(m.logs, m.logs[drop:])
		m.logs = m.logs[:maxLogs]
	}
	m.mu.Unlock()
}

// ClearLogs removes all stored access log entries.
func (m *Monitor) ClearLogs() {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.logs = m.logs[:0]
	m.mu.Unlock()
}

// Snapshot returns traffic history, cumulative stats, and logs (newest first).
func (m *Monitor) Snapshot(maxReturnedLogs int) MonitorSnapshot {
	if m == nil {
		return MonitorSnapshot{}
	}
	m.mu.RLock()
	history := make([]TrafficPoint, len(m.history))
	copy(history, m.history)

	totalLogs := len(m.logs)
	limit := totalLogs
	if maxReturnedLogs > 0 && maxReturnedLogs < limit {
		limit = maxReturnedLogs
	}
	entries := make([]*AccessEntry, limit)
	for i := 0; i < limit; i++ {
		entries[i] = m.logs[totalLogs-1-i]
	}
	m.mu.RUnlock()

	logs := make([]AccessRecord, len(entries))
	for i, entry := range entries {
		logs[i] = entry.snapshot()
	}

	var curUp, curDown int64
	if len(history) > 0 {
		last := history[len(history)-1]
		curUp = last.UpBps
		curDown = last.DownBps
	}

	return MonitorSnapshot{
		TotalUpBytes:    m.totalUpBytes.Load(),
		TotalDownBytes:  m.totalDownBytes.Load(),
		CurrentUpBps:    curUp,
		CurrentDownBps:  curDown,
		ActiveConns:     m.activeConns.Load(),
		TotalRequests:   m.totalRequests.Load(),
		BlockedRequests: m.blockedRequests.Load(),
		History:         history,
		Logs:            logs,
	}
}
