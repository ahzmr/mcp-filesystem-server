package filesystemserver

import (
	"testing"
)

func TestIsToolEnabled(t *testing.T) {
	tests := []struct {
		name       string
		toolName   string
		config     *ToolConfig
		expected   bool
	}{
		{
			name:     "nil config enables all",
			toolName: "read_file",
			config:   nil,
			expected: true,
		},
		{
			name:     "enable all config",
			toolName: "read_file",
			config:   &ToolConfig{EnableAll: true},
			expected: true,
		},
		{
			name:     "exact match",
			toolName: "read_file",
			config:   &ToolConfig{EnabledTools: []string{"read_file", "write_file"}, EnableAll: false},
			expected: true,
		},
		{
			name:     "no match",
			toolName: "delete_file",
			config:   &ToolConfig{EnabledTools: []string{"read_file", "write_file"}, EnableAll: false},
			expected: false,
		},
		{
			name:     "wildcard match - read_*",
			toolName: "read_file",
			config:   &ToolConfig{EnabledTools: []string{"read_*", "write_*"}, EnableAll: false},
			expected: true,
		},
		{
			name:     "wildcard match - read_multiple_files",
			toolName: "read_multiple_files",
			config:   &ToolConfig{EnabledTools: []string{"read_*"}, EnableAll: false},
			expected: true,
		},
		{
			name:     "wildcard no match",
			toolName: "delete_file",
			config:   &ToolConfig{EnabledTools: []string{"read_*", "write_*"}, EnableAll: false},
			expected: false,
		},
		{
			name:     "wildcard match - list_*",
			toolName: "list_directory",
			config:   &ToolConfig{EnabledTools: []string{"list_*"}, EnableAll: false},
			expected: true,
		},
		{
			name:     "wildcard match - list_allowed_directories",
			toolName: "list_allowed_directories",
			config:   &ToolConfig{EnabledTools: []string{"list_*"}, EnableAll: false},
			expected: true,
		},
		{
			name:     "mixed exact and wildcard",
			toolName: "tree",
			config:   &ToolConfig{EnabledTools: []string{"read_*", "tree", "write_*"}, EnableAll: false},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isToolEnabled(tt.toolName, tt.config)
			if result != tt.expected {
				t.Errorf("isToolEnabled(%q, %+v) = %v, want %v", tt.toolName, tt.config, result, tt.expected)
			}
		})
	}
}