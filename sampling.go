package trace

import (
	"fmt"
	"strings"

	. "github.com/infrago/base"
)

type sampleRule struct {
	Name   string
	Sample float64
	Attrs  Map
}

type samplePolicy struct {
	ErrorAlways bool
	HashBy      string
	Rules       []sampleRule
}

func buildSamplePolicy(setting map[string]any) samplePolicy {
	policy := samplePolicy{
		ErrorAlways: true,
		HashBy:      "trace_id",
		Rules:       make([]sampleRule, 0),
	}
	if setting == nil {
		return policy
	}
	if v, ok := setting["error"].(bool); ok {
		policy.ErrorAlways = v
	} else if v, ok := setting["sample_error"].(bool); ok {
		policy.ErrorAlways = v
	}
	if v, ok := setting["key"].(string); ok && strings.TrimSpace(v) != "" {
		policy.HashBy = strings.TrimSpace(strings.ToLower(v))
	} else if v, ok := setting["sample_key"].(string); ok && strings.TrimSpace(v) != "" {
		policy.HashBy = strings.TrimSpace(strings.ToLower(v))
	}

	raw := any(nil)
	if v, ok := setting["rules"]; ok && v != nil {
		raw = v
	} else if v, ok := setting["sample_rules"]; ok && v != nil {
		raw = v
	}
	items, ok := raw.([]any)
	if !ok {
		return policy
	}
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok || m == nil {
			continue
		}
		rule := sampleRule{
			Name:   "",
			Sample: 1,
			Attrs:  Map{},
		}
		if v, ok := m["name"].(string); ok {
			rule.Name = strings.TrimSpace(v)
		}
		if v, ok := parseFloat(m["sample"]); ok {
			rule.Sample = clamp01(v)
		}
		if attrs, ok := m["attrs"].(map[string]any); ok && attrs != nil {
			for k, v := range attrs {
				rule.Attrs[k] = v
			}
		}
		policy.Rules = append(policy.Rules, rule)
	}
	return policy
}

func chooseSampleRatio(span Span, baseRatio float64, policy samplePolicy) float64 {
	if policy.ErrorAlways && (span.Code != 0 || span.Status == StatusFail || span.Status == StatusError) {
		return 1
	}
	ratio := clamp01(baseRatio)
	for _, rule := range policy.Rules {
		if !sampleRuleMatch(rule, span) {
			continue
		}
		ratio = clamp01(rule.Sample)
		break
	}
	return ratio
}

func sampleRuleMatch(rule sampleRule, span Span) bool {
	if rule.Name != "" {
		if strings.HasSuffix(rule.Name, "*") {
			prefix := strings.TrimSuffix(rule.Name, "*")
			if !strings.HasPrefix(span.Name, prefix) {
				return false
			}
		} else if span.Name != rule.Name {
			return false
		}
	}
	if len(rule.Attrs) > 0 {
		for key, expected := range rule.Attrs {
			actual, ok := span.Attributes[key]
			if !ok {
				return false
			}
			if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
				return false
			}
		}
	}
	return true
}

func sampleKey(span Span, hashBy string) string {
	switch hashBy {
	case "span_id":
		if span.SpanId != "" {
			return span.SpanId
		}
	case "trace_span":
		if span.TraceId != "" || span.SpanId != "" {
			return span.TraceId + ":" + span.SpanId
		}
	default:
		if span.TraceId != "" {
			return span.TraceId
		}
	}
	if span.TraceId != "" {
		return span.TraceId
	}
	if span.SpanId != "" {
		return span.SpanId
	}
	return span.Name
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
