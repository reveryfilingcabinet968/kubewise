// Copyright 2026 KubeWise Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package units

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	kilobyte int64 = 1000
	megabyte       = 1000 * kilobyte
	gigabyte       = 1000 * megabyte
	terabyte       = 1000 * gigabyte

	kibibyte int64 = 1024
	mebibyte       = 1024 * kibibyte
	gibibyte       = 1024 * mebibyte
	tebibyte       = 1024 * gibibyte
)

// BytesToGiB converts bytes to gibibytes.
func BytesToGiB(b int64) float64 {
	return float64(b) / float64(gibibyte)
}

// ParseMemory parses a Kubernetes memory string and returns bytes.
// Accepted formats: "512Mi", "2Gi", "1Ti", "500Ki",
// "1000M", "2G", "1T", "500K", "1073741824" (raw bytes).
func ParseMemory(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty memory string")
	}

	// Binary suffixes (IEC): Ki, Mi, Gi, Ti
	for _, pair := range []struct {
		suffix     string
		multiplier int64
	}{
		{"Ti", tebibyte},
		{"Gi", gibibyte},
		{"Mi", mebibyte},
		{"Ki", kibibyte},
	} {
		if raw, ok := strings.CutSuffix(s, pair.suffix); ok {
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid memory value %q: %w", s, err)
			}
			return v * pair.multiplier, nil
		}
	}

	// Decimal suffixes (SI): K, M, G, T
	for _, pair := range []struct {
		suffix     string
		multiplier int64
	}{
		{"T", terabyte},
		{"G", gigabyte},
		{"M", megabyte},
		{"K", kilobyte},
	} {
		if raw, ok := strings.CutSuffix(s, pair.suffix); ok {
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid memory value %q: %w", s, err)
			}
			return v * pair.multiplier, nil
		}
	}

	// Raw bytes
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory value %q: %w", s, err)
	}
	return v, nil
}
