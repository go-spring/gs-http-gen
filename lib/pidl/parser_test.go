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
			input: "/users/profile",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Static, Value: "profile"},
			},
		},
		{
			name:  "colon param",
			input: "/users/:id",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
		},
		{
			name:  "colon wildcard",
			input: "/files/:path*",
			expected: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
		},
		{
			name:  "braced param",
			input: "/users/{id}",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
		},
		{
			name:  "braced wildcard",
			input: "/files/{path...}",
			expected: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
		},
		{
			name:  "mixed params",
			input: "/api/v1/users/{userId}/posts/:postId",
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
			input: "/org/{orgId}/repos/:repoId/branches/{branch...}",
			expected: []Segment{
				{Type: Static, Value: "org"},
				{Type: Param, Value: "orgId"},
				{Type: Static, Value: "repos"},
				{Type: Param, Value: "repoId"},
				{Type: Static, Value: "branches"},
				{Type: Wildcard, Value: "branch"},
			},
		},
		{
			name:     "root path",
			input:    "/",
			expected: []Segment{},
		},
		{
			name:  "single static segment",
			input: "/users",
			expected: []Segment{
				{Type: Static, Value: "users"},
			},
		},
		{
			name:  "multiple static segments",
			input: "/api/v1/users",
			expected: []Segment{
				{Type: Static, Value: "api"},
				{Type: Static, Value: "v1"},
				{Type: Static, Value: "users"},
			},
		},
		{
			name:  "param at beginning",
			input: "/:id/profile",
			expected: []Segment{
				{Type: Param, Value: "id"},
				{Type: Static, Value: "profile"},
			},
		},
		{
			name:  "wildcard at end",
			input: "/files/:path*",
			expected: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
		},
		{
			name:  "braced param with dots",
			input: "/users/{user.id}",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "user.id"},
			},
		},
		{
			name:  "consecutive parameters",
			input: "/users/:userId/posts/:postId/comments/:commentId",
			expected: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "userId"},
				{Type: Static, Value: "posts"},
				{Type: Param, Value: "postId"},
				{Type: Static, Value: "comments"},
				{Type: Param, Value: "commentId"},
			},
		},
		{
			name:     "leading and trailing spaces",
			input:    "  /users/:id  ",
			expected: []Segment{{Type: Static, Value: "users"}, {Type: Param, Value: "id"}},
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

				// Verify that re-parsed results match original
				for i, seg := range reParsedColon {
					if seg.Type != result[i].Type || seg.Value != result[i].Value {
						t.Errorf("re-parsed colon format segment %d mismatch: expected %+v, got %+v", i, result[i], seg)
					}
				}

				for i, seg := range reParsedBrace {
					if seg.Type != result[i].Type || seg.Value != result[i].Value {
						t.Errorf("re-parsed brace format segment %d mismatch: expected %+v, got %+v", i, result[i], seg)
					}
				}
			}

			// Additional test: format with both styles and ensure they parse back correctly
			if len(tt.expected) > 0 {
				colonStyle := Format(tt.expected, Colon)
				braceStyle := Format(tt.expected, Brace)

				colonResult, err := Parse(colonStyle)
				if err != nil {
					t.Errorf("failed to parse colon style %q: %v", colonStyle, err)
				}

				braceResult, err := Parse(braceStyle)
				if err != nil {
					t.Errorf("failed to parse brace style %q: %v", braceStyle, err)
				}

				// Check that both formats produce the same result as expected
				if len(colonResult) != len(tt.expected) {
					t.Errorf("colon format %q produced %d segments, expected %d", colonStyle, len(colonResult), len(tt.expected))
				}

				if len(braceResult) != len(tt.expected) {
					t.Errorf("brace format %q produced %d segments, expected %d", braceStyle, len(braceResult), len(tt.expected))
				}

				for i, seg := range colonResult {
					if seg.Type != tt.expected[i].Type || seg.Value != tt.expected[i].Value {
						t.Errorf("colon format segment %d mismatch: expected %+v, got %+v", i, tt.expected[i], seg)
					}
				}

				for i, seg := range braceResult {
					if seg.Type != tt.expected[i].Type || seg.Value != tt.expected[i].Value {
						t.Errorf("brace format segment %d mismatch: expected %+v, got %+v", i, tt.expected[i], seg)
					}
				}
			}
		})
	}
}

// TestParseErrorCases tests error cases in Parse function
func TestParseErrorCases(t *testing.T) {
	errorTests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid param starting with number",
			input: "/users/:123id",
		},
		{
			name:  "invalid braced param starting with number",
			input: "/users/{123id}",
		},
		{
			name:  "unmatched brace",
			input: "/users/{id",
		},
		{
			name:  "extra wildcard character",
			input: "/users/:id**",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q but got none", tt.input)
			}
		})
	}
}
