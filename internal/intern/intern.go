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

package intern

// StringInterner deduplicates strings to reduce heap allocations.
// In a large cluster, thousands of pods share the same namespace, owner kind,
// label keys, etc. Interning avoids separate heap allocations for each copy.
type StringInterner struct {
	pool map[string]string
}

// New creates a new StringInterner.
func New() *StringInterner {
	return &StringInterner{pool: make(map[string]string, 256)}
}

// Intern returns a deduplicated version of the string.
// If the string was seen before, the previously stored copy is returned.
func (si *StringInterner) Intern(s string) string {
	if existing, ok := si.pool[s]; ok {
		return existing
	}
	si.pool[s] = s
	return s
}

// Len returns the number of unique strings interned.
func (si *StringInterner) Len() int {
	return len(si.pool)
}
