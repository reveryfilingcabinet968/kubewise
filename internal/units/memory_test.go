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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBytesToGiB(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected float64
	}{
		{"zero", 0, 0.0},
		{"1 GiB", 1073741824, 1.0},
		{"2 GiB", 2147483648, 2.0},
		{"512 MiB", 536870912, 0.5},
		{"256 MiB", 268435456, 0.25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, BytesToGiB(tt.input), 1e-9)
		})
	}
}

func TestParseMemory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		// Binary suffixes (IEC)
		{"kibibytes", "512Ki", 524288, false},
		{"mebibytes", "512Mi", 536870912, false},
		{"gibibytes", "2Gi", 2147483648, false},
		{"tebibytes", "1Ti", 1099511627776, false},

		// Decimal suffixes (SI)
		{"kilobytes", "500K", 500000, false},
		{"megabytes", "1000M", 1000000000, false},
		{"gigabytes", "2G", 2000000000, false},
		{"terabytes", "1T", 1000000000000, false},

		// Raw bytes
		{"raw bytes", "1073741824", 1073741824, false},
		{"raw zero", "0", 0, false},

		// With spaces
		{"with spaces", "  512Mi  ", 536870912, false},

		// Errors
		{"empty string", "", 0, true},
		{"invalid", "abc", 0, true},
		{"invalid Mi", "abcMi", 0, true},
		{"invalid G", "abcG", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemory(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
