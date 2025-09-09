/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tidl

import (
	"slices"
)

// GetOneOfAnnotation searches through a slice of annotations and returns
// the first annotation whose key matches any of the provided names.
func GetOneOfAnnotation(arr []Annotation, names ...string) (Annotation, bool) {
	for _, a := range arr {
		if slices.Contains(names, a.Key) {
			return a, true
		}
	}
	return Annotation{}, false
}
