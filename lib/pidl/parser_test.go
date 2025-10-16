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

			if len(result) > 0 {
				// 使用Colon格式化并重新解析
				colonFormatted := Format(result, Colon)
				reParsedColon, err := Parse(colonFormatted)
				if err != nil {
					t.Errorf("failed to re-parse colon formatted path %q: %v", colonFormatted, err)
					return
				}

				// 使用Brace格式化并重新解析
				braceFormatted := Format(result, Brace)
				reParsedBrace, err := Parse(braceFormatted)
				if err != nil {
					t.Errorf("failed to re-parse brace formatted path %q: %v", braceFormatted, err)
					return
				}

				// 验证重新解析的结果与原始结果一致
				if len(reParsedColon) != len(result) || len(reParsedBrace) != len(result) {
					t.Errorf("re-parsed segments count mismatch")
					return
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
		})
	}
}
