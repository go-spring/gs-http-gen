package pidl

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Segment
		hasError bool
	}{
		{
			name:  "empty path",
			input: "",
		},
		{
			name:  "static path",
			input: "users/profile",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Static, Value: "profile"},
			},
		},
		{
			name:  "colon param",
			input: "users/:id",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
		},
		{
			name:  "colon wildcard",
			input: "files/:path*",
			expected: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
		},
		{
			name:  "braced param",
			input: "users/{id}",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
		},
		{
			name:  "braced wildcard",
			input: "files/{path...}",
			expected: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
		},
		{
			name:  "mixed params",
			input: "api/v1/users/{userId}/posts/:postId",
			expected: []Segment{
				{Type: Static, Value: "api"},
				{Type: Static, Value: "v1"},
				{Type: Static, Value: "users"},
				{Type: Param, Value: "userId"},
				{Type: Static, Value: "posts"},
				{Type: Param, Value: "postId"},
			},
		},
		{
			name:  "complex path with multiple param types",
			input: "org/{orgId}/repos/:repoId/branches/{branch...}",
			expected: []Segment{
				{Type: Static, Value: "org"},
				{Type: Param, Value: "orgId"},
				{Type: Static, Value: "repos"},
				{Type: Param, Value: "repoId"},
				{Type: Static, Value: "branches"},
				{Type: Wildcard, Value: "branch"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d segments, got %d", len(tt.expected), len(result))
				return
			}

			for i, seg := range result {
				if seg.Type != tt.expected[i].Type {
					t.Errorf("segment %d: expected type %v, got %v", i, tt.expected[i].Type, seg.Type)
				}
				if seg.Value != tt.expected[i].Value {
					t.Errorf("segment %d: expected value %q, got %q", i, tt.expected[i].Value, seg.Value)
				}
			}

			// Verify using Format function with both styles
			if len(result) > 0 {
				colonFormatted := Format(result, Colon)
				reParsedColon, err := Parse(colonFormatted)
				if err != nil {
					t.Errorf("failed to re-parse colon formatted path %q: %v", colonFormatted, err)
				}

				braceFormatted := Format(result, Brace)
				reParsedBrace, err := Parse(braceFormatted)
				if err != nil {
					t.Errorf("failed to re-parse brace formatted path %q: %v", braceFormatted, err)
				}

				if len(reParsedColon) != len(result) || len(reParsedBrace) != len(result) {
					t.Errorf("re-parsed segments count mismatch")
				}
			}
		})
	}
}
