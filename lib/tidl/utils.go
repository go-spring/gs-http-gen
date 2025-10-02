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

// CapitalizeASCII capitalizes the first ASCII letter of a string.
func CapitalizeASCII(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-'a'+'A') + s[1:]
	}
	return s
}

// GetEnum searches all documents for an enum type with the given name.
func GetEnum(files map[string]Document, name string) (Enum, bool) {
	for _, doc := range files {
		for _, e := range doc.Enums {
			if CapitalizeASCII(e.Name) == name {
				return e, true
			}
		}
	}
	return Enum{}, false
}

// GetType searches all documents for a type with the given name.
func GetType(files map[string]Document, name string) (Type, bool) {
	for _, doc := range files {
		for _, t := range doc.Types {
			if CapitalizeASCII(t.Name) == name {
				return t, true
			}
		}
	}
	return Type{}, false
}

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
