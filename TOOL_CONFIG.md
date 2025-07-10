# 工具配置功能

本文档描述了MCP文件系统服务器的工具配置功能，允许您控制哪些工具可用。

## 概述

MCP文件系统服务器现在支持配置可用的工具列表。您可以：

- 启用所有工具（默认行为）
- 启用特定的工具列表
- 使用通配符模式启用工具组

## 使用方法

### 命令行参数

使用 `--tools` 参数来配置可用的工具：

```bash
# 启用所有工具（默认）
./mcp-filesystem-server /path/to/directory
./mcp-filesystem-server --tools all /path/to/directory

# 启用特定工具
./mcp-filesystem-server --tools read_file,write_file /path/to/directory

# 使用通配符启用工具组
./mcp-filesystem-server --tools 'read_*,list_*' /path/to/directory

# 混合使用精确匹配和通配符
./mcp-filesystem-server --tools 'read_*,tree,delete_file' /path/to/directory
```

### 配置选项

#### 1. 启用所有工具

```bash
--tools all
# 或者省略 --tools 参数
```

这是默认行为，将启用所有可用的工具。

#### 2. 启用特定工具

```bash
--tools read_file,write_file,list_directory
```

使用逗号分隔的工具名称列表。支持的工具包括：

- `read_file` - 读取文件内容
- `write_file` - 写入文件内容
- `list_directory` - 列出目录内容
- `create_directory` - 创建目录
- `copy_file` - 复制文件或目录
- `move_file` - 移动或重命名文件或目录
- `search_files` - 搜索文件和目录
- `get_file_info` - 获取文件或目录信息
- `list_allowed_directories` - 列出允许访问的目录
- `read_multiple_files` - 批量读取多个文件
- `tree` - 获取目录树结构
- `delete_file` - 删除文件或目录
- `modify_file` - 修改文件内容
- `search_within_files` - 在文件内容中搜索文本

#### 3. 使用通配符模式

```bash
--tools 'read_*,write_*'
```

支持的通配符模式：

- `read_*` - 匹配所有以 "read_" 开头的工具（如 `read_file`, `read_multiple_files`）
- `write_*` - 匹配所有以 "write_" 开头的工具（如 `write_file`）
- `list_*` - 匹配所有以 "list_" 开头的工具（如 `list_directory`, `list_allowed_directories`）
- `search_*` - 匹配所有以 "search_" 开头的工具（如 `search_files`, `search_within_files`）

## 示例用例

### 只读访问

```bash
./mcp-filesystem-server --tools 'read_*,list_*,get_file_info,tree' /path/to/directory
```

这将只启用读取和列表相关的工具，不允许修改文件系统。

### 基本文件操作

```bash
./mcp-filesystem-server --tools 'read_file,write_file,list_directory,create_directory' /path/to/directory
```

启用基本的文件读写和目录操作。

### 搜索和分析

```bash
./mcp-filesystem-server --tools 'read_*,search_*,tree,get_file_info' /path/to/directory
```

启用文件读取、搜索和分析相关的工具。

## 安全考虑

工具配置功能可以帮助您：

1. **限制权限** - 只启用必要的工具，减少潜在的安全风险
2. **简化接口** - 为特定用例提供简化的工具集
3. **防止误操作** - 禁用危险的操作（如删除文件）

## 注意事项

- 工具名称区分大小写
- 通配符使用Go的 `filepath.Match` 函数，支持 `*` 和 `?` 通配符
- 空格会被自动忽略，所以 `read_file, write_file` 和 `read_file,write_file` 是等效的
- 如果指定了无效的工具名称，该工具将被忽略（不会报错）