package filesystemserver

import (
	"path/filepath"

	"github.com/mark3labs/mcp-filesystem-server/filesystemserver/handler"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Version = "dev"

// ToolConfig represents the configuration for enabled tools
type ToolConfig struct {
	EnabledTools []string
	EnableAll    bool
}

// isToolEnabled checks if a tool should be enabled based on configuration
func isToolEnabled(toolName string, config *ToolConfig) bool {
	if config == nil || config.EnableAll {
		return true
	}

	for _, pattern := range config.EnabledTools {
		// Check for exact match
		if pattern == toolName {
			return true
		}
		// Check for wildcard match
		if matched, _ := filepath.Match(pattern, toolName); matched {
			return true
		}
	}
	return false
}

func NewFilesystemServer(allowedDirs []string, toolConfig *ToolConfig) (*server.MCPServer, error) {

	h, err := handler.NewFilesystemHandler(allowedDirs)
	if err != nil {
		return nil, err
	}

	s := server.NewMCPServer(
		"secure-filesystem-server",
		Version,
		server.WithResourceCapabilities(true, true),
	)

	// Register resource handlers
	s.AddResource(mcp.NewResource(
		"file://",
		"File System",
		mcp.WithResourceDescription("Access to files and directories on the local file system"),
	), h.HandleReadResource)

	// Register tool handlers
	if isToolEnabled("read_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"read_file",
			mcp.WithDescription("Read the complete contents of a file from the file system."),
			mcp.WithString("path",
				mcp.Description("Path to the file to read"),
				mcp.Required(),
			),
		), h.HandleReadFile)
	}

	if isToolEnabled("write_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"write_file",
			mcp.WithDescription("Create a new file or overwrite an existing file with new content."),
			mcp.WithString("path",
				mcp.Description("Path where to write the file"),
				mcp.Required(),
			),
			mcp.WithString("content",
				mcp.Description("Content to write to the file"),
				mcp.Required(),
			),
		), h.HandleWriteFile)
	}

	if isToolEnabled("list_directory", toolConfig) {
		s.AddTool(mcp.NewTool(
			"list_directory",
			mcp.WithDescription("Get a detailed listing of all files and directories in a specified path."),
			mcp.WithString("path",
				mcp.Description("Path of the directory to list"),
				mcp.Required(),
			),
		), h.HandleListDirectory)
	}

	if isToolEnabled("create_directory", toolConfig) {
		s.AddTool(mcp.NewTool(
			"create_directory",
			mcp.WithDescription("Create a new directory or ensure a directory exists."),
			mcp.WithString("path",
				mcp.Description("Path of the directory to create"),
				mcp.Required(),
			),
		), h.HandleCreateDirectory)
	}

	if isToolEnabled("copy_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"copy_file",
			mcp.WithDescription("Copy files and directories."),
			mcp.WithString("source",
				mcp.Description("Source path of the file or directory"),
				mcp.Required(),
			),
			mcp.WithString("destination",
				mcp.Description("Destination path"),
				mcp.Required(),
			),
		), h.HandleCopyFile)
	}

	if isToolEnabled("move_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"move_file",
			mcp.WithDescription("Move or rename files and directories."),
			mcp.WithString("source",
				mcp.Description("Source path of the file or directory"),
				mcp.Required(),
			),
			mcp.WithString("destination",
				mcp.Description("Destination path"),
				mcp.Required(),
			),
		), h.HandleMoveFile)
	}

	if isToolEnabled("search_files", toolConfig) {
		s.AddTool(mcp.NewTool(
			"search_files",
			mcp.WithDescription("Recursively search for files and directories matching a pattern."),
			mcp.WithString("path",
				mcp.Description("Starting path for the search"),
				mcp.Required(),
			),
			mcp.WithString("pattern",
				mcp.Description("Search pattern to match against file names"),
				mcp.Required(),
			),
		), h.HandleSearchFiles)
	}

	if isToolEnabled("get_file_info", toolConfig) {
		s.AddTool(mcp.NewTool(
			"get_file_info",
			mcp.WithDescription("Retrieve detailed metadata about a file or directory."),
			mcp.WithString("path",
				mcp.Description("Path to the file or directory"),
				mcp.Required(),
			),
		), h.HandleGetFileInfo)
	}

	if isToolEnabled("list_allowed_directories", toolConfig) {
		s.AddTool(mcp.NewTool(
			"list_allowed_directories",
			mcp.WithDescription("Returns the list of directories that this server is allowed to access."),
		), h.HandleListAllowedDirectories)
	}

	if isToolEnabled("read_multiple_files", toolConfig) {
		s.AddTool(mcp.NewTool(
			"read_multiple_files",
			mcp.WithDescription("Read the contents of multiple files in a single operation."),
			mcp.WithArray("paths",
				mcp.Description("List of file paths to read"),
				mcp.Required(),
				mcp.Items(map[string]any{"type": "string"}),
			),
		), h.HandleReadMultipleFiles)
	}

	if isToolEnabled("tree", toolConfig) {
		s.AddTool(mcp.NewTool(
			"tree",
			mcp.WithDescription("Returns a hierarchical JSON representation of a directory structure."),
			mcp.WithString("path",
				mcp.Description("Path of the directory to traverse"),
				mcp.Required(),
			),
			mcp.WithNumber("depth",
				mcp.Description("Maximum depth to traverse (default: 3)"),
			),
			mcp.WithBoolean("follow_symlinks",
				mcp.Description("Whether to follow symbolic links (default: false)"),
			),
		), h.HandleTree)
	}

	if isToolEnabled("delete_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"delete_file",
			mcp.WithDescription("Delete a file or directory from the file system."),
			mcp.WithString("path",
				mcp.Description("Path to the file or directory to delete"),
				mcp.Required(),
			),
			mcp.WithBoolean("recursive",
				mcp.Description("Whether to recursively delete directories (default: false)"),
			),
		), h.HandleDeleteFile)
	}

	if isToolEnabled("modify_file", toolConfig) {
		s.AddTool(mcp.NewTool(
			"modify_file",
			mcp.WithDescription("Update file by finding and replacing text. Provides a simple pattern matching interface without needing exact character positions."),
			mcp.WithString("path",
				mcp.Description("Path to the file to modify"),
				mcp.Required(),
			),
			mcp.WithString("find",
				mcp.Description("Text to search for (exact match or regex pattern)"),
				mcp.Required(),
			),
			mcp.WithString("replace",
				mcp.Description("Text to replace with"),
				mcp.Required(),
			),
			mcp.WithBoolean("all_occurrences",
				mcp.Description("Replace all occurrences of the matching text (default: true)"),
			),
			mcp.WithBoolean("regex",
				mcp.Description("Treat the find pattern as a regular expression (default: false)"),
			),
		), h.HandleModifyFile)
	}

	if isToolEnabled("search_within_files", toolConfig) {
		s.AddTool(mcp.NewTool(
			"search_within_files",
			mcp.WithDescription("Search for text within file contents. Unlike search_files which only searches file names, this tool scans the actual contents of text files for matching substrings. Binary files are automatically excluded from the search. Reports file paths and line numbers where matches are found."),
			mcp.WithString("path",
				mcp.Description("Starting path for the search (must be a directory)"),
				mcp.Required(),
			),
			mcp.WithString("substring",
				mcp.Description("Text to search for within file contents"),
				mcp.Required(),
			),
			mcp.WithNumber("depth",
				mcp.Description("Maximum directory depth to search (default: unlimited)"),
			),
			mcp.WithNumber("max_results",
				mcp.Description("Maximum number of results to return (default: 1000)"),
			),
		), h.HandleSearchWithinFiles)
	}

	return s, nil
}
