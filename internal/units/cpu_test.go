//    Copyright 2026 KubeWise Authors
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMillicoresToCores(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected float64
	}{
		{"zero", 0, 0.0},
		{"one core", 1000, 1.0},
		{"half core", 500, 0.5},
		{"250m", 250, 0.25},
		{"4 cores", 4000, 4.0},
		{"1m", 1, 0.001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, MillicoresToCores(tt.input), 1e-9)
		})
	}
}

func TestParseCPU(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"millicores", "500m", 500, false},
		{"millicores zero", "0m", 0, false},
		{"millicores large", "4000m", 4000, false},
		{"whole cores", "2", 2000, false},
		{"fractional cores", "0.5", 500, false},
		{"one core", "1", 1000, false},
		{"quarter core", "0.25", 250, false},
		{"with spaces", "  500m  ", 500, false},
		{"empty string", "", 0, true},
		{"invalid", "abc", 0, true},
		{"invalid millicore", "abcm", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCPU(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
