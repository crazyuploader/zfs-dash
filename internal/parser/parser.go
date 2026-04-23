// Package parser implements a minimal Prometheus text-format exposition parser.
package parser

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Sample is one metric observation.
type Sample struct {
	Name   string
	Labels map[string]string
	Value  float64
}

// Parse reads Prometheus text format from r and returns all samples.
func Parse(r io.Reader) ([]Sample, error) {
	var samples []Sample
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		s, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		samples = append(samples, s)
	}
	return samples, scanner.Err()
}

func parseLine(line string) (Sample, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Sample{}, io.ErrUnexpectedEOF
	}
	val, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Sample{}, err
	}

	raw := parts[0]
	name := raw
	labels := map[string]string{}

	if idx := strings.IndexByte(raw, '{'); idx >= 0 {
		name = raw[:idx]
		inner := raw[idx+1:]
		if end := strings.LastIndexByte(inner, '}'); end >= 0 {
			inner = inner[:end]
		}
		for _, pair := range splitLabels(inner) {
			pair = strings.TrimSpace(pair)
			if eq := strings.IndexByte(pair, '='); eq >= 0 {
				k := strings.TrimSpace(pair[:eq])
				v := strings.Trim(strings.TrimSpace(pair[eq+1:]), `"`)
				labels[k] = v
			}
		}
	}
	return Sample{Name: name, Labels: labels, Value: val}, nil
}

// splitLabels splits a label string by commas while respecting quoted values.
func splitLabels(s string) []string {
	var out []string
	depth, start := 0, 0
	for i, ch := range s {
		switch ch {
		case '"':
			depth ^= 1
		case ',':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
