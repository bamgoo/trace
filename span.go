package trace

import (
	"fmt"
	"hash/fnv"
	"sync/atomic"
	"time"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
)

type Span struct {
	Time              time.Time `json:"time"`
	TraceId           string    `json:"trace_id,omitempty"`
	SpanId            string    `json:"span_id,omitempty"`
	ParentSpanId      string    `json:"parent_span_id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Kind              string    `json:"kind,omitempty"`
	ServiceName       string    `json:"service_name,omitempty"`
	Target            string    `json:"target,omitempty"`
	Status            string    `json:"status,omitempty"`
	StatusCode        string    `json:"status_code,omitempty"`
	StatusMessage     string    `json:"status_message,omitempty"`
	DurationMs        int64     `json:"duration_ms"`
	StartMs           int64     `json:"start_ms"`
	EndMs             int64     `json:"end_ms"`
	StartTimeUnixNano int64     `json:"start_time_unix_nano"`
	EndTimeUnixNano   int64     `json:"end_time_unix_nano"`
	Timestamp         int64     `json:"timestamp"`
	Attributes        Map       `json:"attributes,omitempty"`
	Resource          Map       `json:"resource,omitempty"`
}

type Handle struct {
	meta  *bamgoo.Meta
	span  Span
	start time.Time
}

var spanSeq atomic.Uint64

func Begin(meta *bamgoo.Meta, name string, attrs ...Map) *Handle {
	if meta == nil {
		meta = bamgoo.NewMeta()
	}
	traceId := meta.TraceId()
	if traceId == "" {
		traceId = nextId("tr")
		meta.TraceId(traceId)
	}
	parent := meta.SpanId()
	spanId := nextId("sp")
	prevParent := meta.ParentSpanId()
	meta.ParentSpanId(parent)
	meta.SpanId(spanId)
	meta.PushSpanFrame(parent, prevParent)

	span := Span{
		TraceId:           traceId,
		SpanId:            spanId,
		ParentSpanId:      parent,
		Name:              name,
		Status:            StatusOK,
		StatusCode:        "STATUS_CODE_OK",
		StartMs:           time.Now().UnixMilli(),
		StartTimeUnixNano: time.Now().UnixNano(),
		Attributes:        Map{},
		Resource:          Map{},
	}
	identity := bamgoo.Identity()
	span.Resource["bamgoo.project"] = identity.Project
	span.Resource["bamgoo.profile"] = identity.Profile
	span.Resource["bamgoo.node"] = identity.Node
	if len(attrs) > 0 && attrs[0] != nil {
		for k, v := range attrs[0] {
			span.Attributes[k] = v
		}
		if v, ok := attrs[0]["kind"].(string); ok {
			span.Kind = v
		}
		if v, ok := attrs[0]["service"].(string); ok {
			span.ServiceName = v
			span.Resource["service.name"] = v
		}
		if v, ok := attrs[0]["target"].(string); ok {
			span.Target = v
		}
		if v, ok := attrs[0]["status_code"].(string); ok && v != "" {
			span.StatusCode = v
		}
		if v, ok := attrs[0]["status_message"].(string); ok && v != "" {
			span.StatusMessage = v
		}
	}

	return &Handle{
		meta:  meta,
		span:  span,
		start: time.Now(),
	}
}

func (h *Handle) End(errs ...error) {
	if h == nil {
		return
	}
	now := time.Now()
	h.span.Time = now
	h.span.EndMs = now.UnixMilli()
	h.span.EndTimeUnixNano = now.UnixNano()
	h.span.DurationMs = time.Since(h.start).Milliseconds()
	h.span.Timestamp = now.UnixMilli()
	if len(errs) > 0 && errs[0] != nil {
		h.span.Status = StatusError
		h.span.StatusCode = "STATUS_CODE_ERROR"
		h.span.StatusMessage = errs[0].Error()
	}
	Write(h.span)
	if h.meta != nil {
		if prevSpanId, prevParent, ok := h.meta.PopSpanFrame(); ok {
			h.meta.SpanId(prevSpanId)
			h.meta.ParentSpanId(prevParent)
		}
	}
}

func Emit(meta *bamgoo.Meta, name string, status string, attrs ...Map) {
	h := Begin(meta, name, attrs...)
	if h == nil {
		return
	}
	if status != "" {
		h.span.Status = status
		if status == StatusError {
			h.span.StatusCode = "STATUS_CODE_ERROR"
		}
	}
	h.End()
}

func nextId(prefix string) string {
	n := spanSeq.Add(1)
	return fmt.Sprintf("%s%x%x", prefix, time.Now().UnixNano(), n)
}

func hash01(s string) float64 {
	if s == "" {
		return 1
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	v := h.Sum64()
	return float64(v%1000000) / 1000000.0
}
