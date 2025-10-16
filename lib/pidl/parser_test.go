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
			hasError: true,
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
			hasError: true,
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
		{
			name:     "invalid param starting with number",
			input:    "/users/:123id",
			hasError: true,
		},
		{
			name:     "invalid braced param starting with number",
			input:    "/users/{123id}",
			hasError: true,
		},
		{
			name:     "unmatched brace",
			input:    "/users/{id",
			hasError: true,
		},
		{
			name:     "extra wildcard character",
			input:    "/users/:id**",
			hasError: true,
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
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    []Segment
		expected string
		style    SegmentStyle
	}{
		{
			name: "static path with colon style",
			input: []Segment{
				{Type: Static, Value: "users"},
				{Type: Static, Value: "profile"},
			},
			expected: "/users/profile",
			style:    Colon,
		},
		{
			name: "param with colon style",
			input: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
			expected: "/users/:id",
			style:    Colon,
		},
		{
			name: "wildcard with colon style",
			input: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
			expected: "/files/:path*",
			style:    Colon,
		},
		{
			name: "param with brace style",
			input: []Segment{
				{Type: Static, Value: "users"},
				{Type: Param, Value: "id"},
			},
			expected: "/users/{id}",
			style:    Brace,
		},
		{
			name: "wildcard with brace style",
			input: []Segment{
				{Type: Static, Value: "files"},
				{Type: Wildcard, Value: "path"},
			},
			expected: "/files/{path...}",
			style:    Brace,
		},
		{
			name: "complex path with colon style",
			input: []Segment{
				{
					Type:  Static,
					Value: "api",
				},
				{Type: Static, Value: "v1"},
				{Type: Static, Value: "users"},
				{Type: Param, Value: "userId"},
				{Type: Static, Value: "posts"},
				{Type: Param, Value: "postId"},
			},
			expected: "/api/v1/users/:userId/posts/:postId",
			style:    Colon,
		},
		{
			name: "complex path with brace style",
			input: []Segment{
				{Type: Static, Value: "api"},
				{Type: Static, Value: "v1"},
				{Type: Static, Value: "users"},
				{Type: Param, Value: "userId"},
				{Type: Static, Value: "posts"},
				{Type: Param, Value: "postId"},
			},
			expected: "/api/v1/users/{userId}/posts/{postId}",
			style:    Brace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.input, tt.style)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
