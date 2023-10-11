package tracer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
)

type Span struct {
	ID           string
	TraceID      string
	ParentSpanID string
	StartTime    time.Time
	EndTime      time.Time
}

type Trace struct {
	ID        string
	Spans     map[string]*Span
	StartTime time.Time
	EndTime   time.Time
	mu        sync.Mutex
}

// NewTrace starts a new trace.
func NewTrace() *Trace {
	return &Trace{
		ID:        generateTraceID(),
		Spans:     make(map[string]*Span),
		StartTime: time.Now(),
	}
}

// StartSpan starts a new root span.
//
// For SubSpan, use StartSubSpan instead.
func (t *Trace) StartSpan() *Span {
	t.mu.Lock()
	defer t.mu.Unlock()

	span := &Span{
		TraceID:      t.ID,
		ID:           generateSpanID(),
		ParentSpanID: "",
		StartTime:    time.Now(),
	}

	t.Spans[span.ID] = span
	return span
}

// StartSubSpan starts a new span as a child of the given parent span.
func (t *Trace) StartSubSpan(parentSpan *Span) (*Span, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, s := range t.Spans {
		if s.ID == parentSpan.ID {
			parentSpan = s
			break
		}
	}

	if parentSpan == nil {
		return nil, fmt.Errorf("parent span with ID %s not found", parentSpan.ID)
	}

	subSpan := &Span{
		TraceID:      t.ID,
		ID:           generateSpanID(),
		ParentSpanID: parentSpan.ID,
		StartTime:    time.Now(),
	}

	t.Spans[subSpan.ID] = subSpan
	return subSpan, nil
}

func generateTraceID() string {
	return uuid.New().String()
}

func generateSpanID() string {
	return shortuuid.New()
}

// EndSpan ends a Span.
func (t *Trace) EndSpan(span *Span) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, exists := t.Spans[span.ID]
	if !exists {
		return errors.New("span does not exist")
	}

	if span.EndTime.IsZero() {
		span.EndTime = time.Now()
	} else {
		return errors.New("span has already ended")
	}

	return nil
}

// EndTrace ends a Trace.
func (t *Trace) EndTrace() (*Trace, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.EndTime.IsZero() {
		t.EndTime = time.Now()
	} else {
		return nil, errors.New("span has already ended")
	}

	return t, nil
}

// Must returns Trace if err is nil and panics otherwise.
func Must(t *Trace, err error) *Trace {
	return t
}
