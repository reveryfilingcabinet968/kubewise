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

// MillicoresToCores converts CPU millicores to cores.
func MillicoresToCores(m int64) float64 {
	return float64(m) / 1000.0
}

// ParseCPU parses a Kubernetes CPU string and returns millicores.
// Accepted formats: "500m" (millicores), "2" (cores), "0.5" (cores).
func ParseCPU(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty CPU string")
	}

	if raw, ok := strings.CutSuffix(s, "m"); ok {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid CPU millicore value %q: %w", s, err)
		}
		return v, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid CPU value %q: %w", s, err)
	}
	return int64(v * 1000), nil
}
