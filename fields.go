package trace

import (
	"fmt"
	"strings"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
)

// SpanValues returns all known trace fields (standard + compat aliases).
func SpanValues(span Span, instance, flag string) map[string]Any {
	identity := bamgoo.Identity()
	project := identity.Project
	profile := identity.Profile
	node := identity.Node
	if span.Resource != nil {
		if v, ok := span.Resource["bamgoo.project"].(string); ok && v != "" {
			project = v
		}
		if v, ok := span.Resource["bamgoo.profile"].(string); ok && v != "" {
			profile = v
		}
		if v, ok := span.Resource["bamgoo.node"].(string); ok && v != "" {
			node = v
		}
	}

	values := map[string]Any{
		"time":                 span.Time.Format("2006-01-02 15:04:05.000"),
		"trace_id":             span.TraceId,
		"span_id":              span.SpanId,
		"parent_span_id":       span.ParentSpanId,
		"name":                 span.Name,
		"kind":                 span.Kind,
		"service_name":         span.ServiceName,
		"target":               span.Target,
		"status":               span.Status,
		"status_code":          span.StatusCode,
		"status_message":       span.StatusMessage,
		"duration_ms":          span.DurationMs,
		"start_ms":             span.StartMs,
		"end_ms":               span.EndMs,
		"start_time_unix_nano": span.StartTimeUnixNano,
		"end_time_unix_nano":   span.EndTimeUnixNano,
		"timestamp":            span.Timestamp,
		"attributes":           span.Attributes,
		"resource":             span.Resource,
		"project":              project,
		"profile":              profile,
		"node":                 node,
		"flag":                 flag,
		"ts":                   span.Time,
		// Compat aliases
		"traceId":      span.TraceId,
		"spanId":       span.SpanId,
		"parentSpanId": span.ParentSpanId,
		"service":      span.ServiceName,
		"error":        span.StatusMessage,
		"cost_ms":      span.DurationMs,
		"costMs":       span.DurationMs,
		"startMs":      span.StartMs,
		"endMs":        span.EndMs,
		"attrs":        span.Attributes,
		"parent_id":    span.ParentSpanId,
		"instance":     instance,
	}
	return values
}

// ResolveFields parses fields config into source->target mapping.
// Supports:
//   - []string / []any: ["trace_id","span_id"]
//   - map[string]any: { trace_id = "tid", span_id = "sid" }
func ResolveFields(raw Any, defaults map[string]string) map[string]string {
	out := cloneFieldMap(defaults)
	if raw == nil {
		return out
	}
	switch vv := raw.(type) {
	case []string:
		out = map[string]string{}
		for _, source := range vv {
			source = strings.TrimSpace(source)
			if source == "" {
				continue
			}
			out[source] = source
		}
	case []any:
		out = map[string]string{}
		for _, item := range vv {
			source := strings.TrimSpace(fmt.Sprintf("%v", item))
			if source == "" {
				continue
			}
			out[source] = source
		}
	case map[string]any:
		out = map[string]string{}
		for source, targetAny := range vv {
			source = strings.TrimSpace(source)
			target := strings.TrimSpace(fmt.Sprintf("%v", targetAny))
			if source == "" || target == "" {
				continue
			}
			out[source] = target
		}
	}

	if len(out) == 0 {
		return cloneFieldMap(defaults)
	}
	return out
}

func cloneFieldMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
