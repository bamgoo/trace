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

	attrs := sanitizeTraceAttrs(span.Attributes)
	values := map[string]Any{
		"time":           span.Time.UnixNano(),
		"start":          span.Start,
		"end":            span.End,
		"cost":           span.Cost,
		"trace_id":       span.TraceId,
		"span_id":        span.SpanId,
		"parent_span_id": span.ParentSpanId,
		"step":           span.Name,
		"kind":           span.Kind,
		"entry":          span.Target,
		"status":         span.Status,
		"code":           span.Code,
		"result":         span.Result,
		"attributes":     attrs,
		"resource":       span.Resource,
		"project":        project,
		"profile":        profile,
		"node":           node,
		"flag":           flag,
		"ts":             span.Time,
		// Compat aliases
		"name":         span.Name,
		"span":         span.Name,
		"traceId":      span.TraceId,
		"spanId":       span.SpanId,
		"parentSpanId": span.ParentSpanId,
		"service":      span.ServiceName,
		"error":        span.Result,
		"costMs":       span.Cost / int64(1e6),
		"startMs":      span.Start / int64(1e6),
		"endMs":        span.End / int64(1e6),
		"attrs":        attrs,
		"parent_id":    span.ParentSpanId,
		"instance":     instance,
	}
	if v, ok := stringAttr(span.Attributes, "step"); ok && v != "" {
		values["step"] = v
		values["name"] = v
		values["span"] = v
	} else if v, ok := stringAttr(span.Attributes, "name"); ok && v != "" {
		values["step"] = v
		values["name"] = v
		values["span"] = v
	}
	if v, ok := stringAttr(span.Attributes, "entry"); ok && v != "" {
		values["entry"] = v
	}
	if v, ok := stringAttr(span.Attributes, "target"); ok && v != "" && values["entry"] == "" {
		values["entry"] = v
	}
	if values["step"] == "" {
		values["step"] = "internal"
		values["name"] = "internal"
		values["span"] = "internal"
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

func stringAttr(m Map, key string) (string, bool) {
	if m == nil || key == "" {
		return "", false
	}
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func sanitizeTraceAttrs(in Map) Map {
	if in == nil {
		return Map{}
	}
	out := Map{}
	for k, v := range in {
		switch k {
		case "kind", "service", "entry", "step", "status", "code", "result", "target":
			continue
		}
		out[k] = v
	}
	return out
}
