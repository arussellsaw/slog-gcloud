package sloggcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/monzo/slog"
)

func NewLogger(ctx context.Context, name string) (*StackDriverLogger, error) {
	l, err := logging.NewClient(ctx, ProjectID)
	if err != nil {
		return nil, err
	}

	return &StackDriverLogger{
		logger: l.Logger(name),
	}, nil
}

// StackDriverLogger is an implementation of monzo/slog.Logger
// that emits stackdriver compatible events
type StackDriverLogger struct {
	mu     sync.Mutex
	buffer []slog.Event
	logger *logging.Logger
}

func (l *StackDriverLogger) Log(evs ...slog.Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range evs {
		l.logger.Log(logging.Entry{
			Timestamp: e.Timestamp,
			Labels:    allLabels(e),
			Trace:     Trace(e.Context),
			Payload:   e.Message,
			Severity:  logging.ParseSeverity(e.Severity.String()),
		})
	}
}

func (l *StackDriverLogger) Flush() error {
	return nil
}

// Entry ...
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	Params map[string]string `json:"params,omitempty"`
}

// String renders an entry structure to the JSON format expected by Stackdriver.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		fmt.Println("json.Marshal:", err)
	}
	return string(out)
}

func NewEntry(e slog.Event) Entry {
	metadata := make(map[string]string)
	for k, v := range e.Metadata {
		metadata[k] = fmt.Sprint(v)
	}
	return Entry{
		Trace:    Trace(e.Context),
		Severity: e.Severity.String(),
		Message:  e.Message,
		Params:   metadata,
	}
}

func allLabels(e slog.Event) map[string]string {
	out := make(map[string]string)
	for k, v := range e.Labels {
		out[k] = v
	}
	for k, v := range e.Metadata {
		out[k] = fmt.Sprint(v)
	}
	return out
}
