// Package gitconfig
// @description: KnownKey 内置 Git 配置键定义库
package gitconfig

// KnownKey 内置键定义
type KnownKey struct {
	Section    string   `json:"section"`
	Key        string   `json:"key"`
	Type       string   `json:"type"`       // "string" | "bool" | "int" | "enum"
	Default    string   `json:"default"`
	EnumValues []string `json:"enumValues"` // type=enum 时有效
	DescZh     string   `json:"descZh"`
	DescEn     string   `json:"descEn"`
}

var knownKeys = []KnownKey{
	// ========== core ==========
	{Section: "core", Key: "repositoryformatversion", Type: "int", Default: "0", DescZh: "仓库格式版本号，Git 内部使用", DescEn: "Repository format version, used internally by Git"},
	{Section: "core", Key: "filemode", Type: "bool", Default: "true", DescZh: "是否跟踪文件权限变更（chmod），Windows 上通常为 false", DescEn: "Track file permission changes (chmod), usually false on Windows"},
	{Section: "core", Key: "bare", Type: "bool", Default: "false", DescZh: "是否为裸仓库（无工作区，仅用于服务端）", DescEn: "Whether this is a bare repository (no working tree, server-side only)"},
	{Section: "core", Key: "logallrefupdates", Type: "bool", Default: "true", DescZh: "是否记录所有引用更新到 reflog，用于恢复误操作", DescEn: "Log all ref updates to reflog for recovery"},
	{Section: "core", Key: "autocrlf", Type: "enum", Default: "false", EnumValues: []string{"true", "false", "input"}, DescZh: "自动换行符转换：true=提交LF检出CRLF，input=仅提交时转LF，false=不转换", DescEn: "Auto CRLF conversion: true=commit LF checkout CRLF, input=commit only, false=no conversion"},
	{Section: "core", Key: "eol", Type: "enum", Default: "", EnumValues: []string{"lf", "crlf", "native"}, DescZh: "工作区文件的换行符格式：lf=Unix，crlf=Windows，native=系统默认", DescEn: "Line ending format: lf=Unix, crlf=Windows, native=system default"},
	{Section: "core", Key: "ignorecase", Type: "bool", Default: "false", DescZh: "是否忽略文件名大小写（macOS/Windows 默认 true）", DescEn: "Ignore filename case (default true on macOS/Windows)"},
	{Section: "core", Key: "editor", Type: "string", Default: "", DescZh: "Git 使用的默认文本编辑器（如 vim、nano、code --wait）", DescEn: "Default text editor (e.g. vim, nano, code --wait)"},
	{Section: "core", Key: "whitespace", Type: "string", Default: "", DescZh: "空白字符检测规则（如 trailing-space、space-before-tab）", DescEn: "Whitespace error detection rules"},
	{Section: "core", Key: "excludesfile", Type: "string", Default: "", DescZh: "全局 .gitignore 文件路径（如 ~/.gitignore_global）", DescEn: "Path to global gitignore file"},
	{Section: "core", Key: "quotepath", Type: "bool", Default: "true", DescZh: "是否对非 ASCII 字符路径加引号转义", DescEn: "Quote non-ASCII characters in paths"},
	{Section: "core", Key: "compression", Type: "int", Default: "-1", DescZh: "压缩级别（-1=默认，0=不压缩，9=最高压缩）", DescEn: "Compression level (-1=default, 0=none, 9=max)"},

	// ========== user ==========
	{Section: "user", Key: "name", Type: "string", Default: "", DescZh: "提交时使用的用户名（显示在 git log 中）", DescEn: "Name used in commits (shown in git log)"},
	{Section: "user", Key: "email", Type: "string", Default: "", DescZh: "提交时使用的邮箱地址（用于识别作者）", DescEn: "Email address used in commits (identifies author)"},
	{Section: "user", Key: "signingkey", Type: "string", Default: "", DescZh: "GPG 签名密钥 ID（用于签名提交和标签）", DescEn: "GPG signing key ID (for signed commits and tags)"},

	// ========== remote ==========
	{Section: "remote", Key: "url", Type: "string", Default: "", DescZh: "远程仓库地址（支持 https://、git://、ssh://）", DescEn: "Remote repository URL (supports https://, git://, ssh://)"},
	{Section: "remote", Key: "fetch", Type: "string", Default: "+refs/heads/*:refs/remotes/origin/*", DescZh: "fetch 时的 refspec 映射规则（定义本地/远程分支对应关系）", DescEn: "Refspec mapping for fetch (defines local/remote branch mapping)"},
	{Section: "remote", Key: "pushurl", Type: "string", Default: "", DescZh: "push 时使用的独立 URL（覆盖 url，用于读写分离）", DescEn: "Separate URL for push (overrides url, for read/write separation)"},
	{Section: "remote", Key: "mirror", Type: "bool", Default: "false", DescZh: "是否以镜像模式推送（推送所有引用，危险操作）", DescEn: "Push in mirror mode (pushes all refs, dangerous)"},
	{Section: "remote", Key: "tagopt", Type: "enum", Default: "", EnumValues: []string{"--tags", "--no-tags"}, DescZh: "fetch 时是否自动获取标签", DescEn: "Auto-fetch tags on fetch"},

	// ========== branch ==========
	{Section: "branch", Key: "remote", Type: "string", Default: "origin", DescZh: "该分支跟踪的远程名称（通常为 origin）", DescEn: "Remote tracked by this branch (usually origin)"},
	{Section: "branch", Key: "merge", Type: "string", Default: "", DescZh: "该分支跟踪的远程引用（如 refs/heads/main）", DescEn: "Remote ref tracked by this branch (e.g. refs/heads/main)"},
	{Section: "branch", Key: "rebase", Type: "bool", Default: "false", DescZh: "pull 时是否默认使用 rebase 而非 merge", DescEn: "Use rebase instead of merge on pull"},

	// ========== pull ==========
	{Section: "pull", Key: "rebase", Type: "enum", Default: "false", EnumValues: []string{"true", "false", "merges", "interactive"}, DescZh: "pull 时的 rebase 策略：true=变基，merges=保留合并提交，interactive=交互式", DescEn: "Rebase strategy on pull: true=rebase, merges=preserve merges, interactive=interactive"},
	{Section: "pull", Key: "ff", Type: "enum", Default: "", EnumValues: []string{"true", "false", "only"}, DescZh: "pull 时的 fast-forward 策略：only=仅快进，false=总是创建合并提交", DescEn: "Fast-forward strategy: only=fast-forward only, false=always create merge commit"},

	// ========== push ==========
	{Section: "push", Key: "default", Type: "enum", Default: "simple", EnumValues: []string{"nothing", "current", "upstream", "tracking", "simple", "matching"}, DescZh: "push 时的默认行为：simple=推送当前分支到同名远程分支，current=推送到同名分支", DescEn: "Default push behavior: simple=push current to same-named remote, current=push to same name"},
	{Section: "push", Key: "followtags", Type: "bool", Default: "false", DescZh: "push 时是否自动推送指向当前提交的标签", DescEn: "Automatically push tags pointing to pushed commits"},
	{Section: "push", Key: "autoSetupRemote", Type: "bool", Default: "false", DescZh: "首次 push 时是否自动设置上游分支", DescEn: "Auto-setup upstream branch on first push"},

	// ========== merge ==========
	{Section: "merge", Key: "tool", Type: "string", Default: "", DescZh: "合并冲突时使用的工具（如 vimdiff、meld、vscode）", DescEn: "Tool for merge conflicts (e.g. vimdiff, meld, vscode)"},
	{Section: "merge", Key: "conflictstyle", Type: "enum", Default: "merge", EnumValues: []string{"merge", "diff3", "zdiff3"}, DescZh: "冲突标记的显示风格：diff3=显示共同祖先，zdiff3=压缩相同行", DescEn: "Conflict marker style: diff3=show common ancestor, zdiff3=compress same lines"},
	{Section: "merge", Key: "ff", Type: "enum", Default: "", EnumValues: []string{"true", "false", "only"}, DescZh: "merge 时的 fast-forward 策略", DescEn: "Fast-forward strategy on merge"},

	// ========== rebase ==========
	{Section: "rebase", Key: "autosquash", Type: "bool", Default: "false", DescZh: "是否自动整理 fixup!/squash! 提交（交互式 rebase 时）", DescEn: "Auto-squash fixup!/squash! commits in interactive rebase"},
	{Section: "rebase", Key: "autostash", Type: "bool", Default: "false", DescZh: "rebase 前是否自动 stash 工作区变更并在完成后恢复", DescEn: "Auto-stash changes before rebase and restore after"},

	// ========== diff ==========
	{Section: "diff", Key: "tool", Type: "string", Default: "", DescZh: "diff 时使用的外部工具（如 vimdiff、meld、beyond compare）", DescEn: "External diff tool (e.g. vimdiff, meld, beyond compare)"},
	{Section: "diff", Key: "algorithm", Type: "enum", Default: "myers", EnumValues: []string{"myers", "minimal", "patience", "histogram"}, DescZh: "diff 算法：histogram=更快更准确（推荐），patience=更易读", DescEn: "Diff algorithm: histogram=faster and accurate (recommended), patience=more readable"},
	{Section: "diff", Key: "renames", Type: "bool", Default: "true", DescZh: "是否检测文件重命名", DescEn: "Detect file renames"},

	// ========== http / https ==========
	{Section: "http", Key: "proxy", Type: "string", Default: "", DescZh: "HTTP 代理地址（如 http://127.0.0.1:7890）", DescEn: "HTTP proxy URL (e.g. http://127.0.0.1:7890)"},
	{Section: "http", Key: "sslverify", Type: "bool", Default: "true", DescZh: "是否验证 SSL 证书（自签名证书可设为 false）", DescEn: "Verify SSL certificates (set false for self-signed certs)"},
	{Section: "http", Key: "sslcainfo", Type: "string", Default: "", DescZh: "自定义 CA 证书文件路径（用于企业内网）", DescEn: "Path to custom CA certificate file (for corporate networks)"},
	{Section: "http", Key: "postbuffer", Type: "int", Default: "1048576", DescZh: "POST 缓冲区大小（字节），推送大文件时可增大", DescEn: "POST buffer size in bytes, increase for large pushes"},
	{Section: "https", Key: "proxy", Type: "string", Default: "", DescZh: "HTTPS 代理地址", DescEn: "HTTPS proxy URL"},

	// ========== credential ==========
	{Section: "credential", Key: "helper", Type: "string", Default: "", DescZh: "凭据存储助手：store=明文存储，cache=内存缓存，osxkeychain/wincred/libsecret=系统钥匙串", DescEn: "Credential helper: store=plaintext, cache=memory, osxkeychain/wincred/libsecret=system keychain"},

	// ========== fetch ==========
	{Section: "fetch", Key: "prune", Type: "bool", Default: "false", DescZh: "fetch 时是否自动删除远程已删除的本地引用", DescEn: "Auto-prune deleted remote refs on fetch"},
	{Section: "fetch", Key: "recurseSubmodules", Type: "enum", Default: "on-demand", EnumValues: []string{"true", "false", "on-demand"}, DescZh: "是否递归 fetch 子模块", DescEn: "Recursively fetch submodules"},

	// ========== alias ==========
	{Section: "alias", Key: "st", Type: "string", Default: "status", DescZh: "自定义命令别名（示例：st = status）", DescEn: "Custom command alias (example: st = status)"},
	{Section: "alias", Key: "co", Type: "string", Default: "checkout", DescZh: "自定义命令别名（示例：co = checkout）", DescEn: "Custom command alias (example: co = checkout)"},
	{Section: "alias", Key: "br", Type: "string", Default: "branch", DescZh: "自定义命令别名（示例：br = branch）", DescEn: "Custom command alias (example: br = branch)"},
	{Section: "alias", Key: "ci", Type: "string", Default: "commit", DescZh: "自定义命令别名（示例：ci = commit）", DescEn: "Custom command alias (example: ci = commit)"},

	// ========== submodule ==========
	{Section: "submodule", Key: "active", Type: "bool", Default: "true", DescZh: "该子模块是否激活", DescEn: "Whether this submodule is active"},
	{Section: "submodule", Key: "url", Type: "string", Default: "", DescZh: "子模块的远程仓库地址", DescEn: "Remote URL of the submodule"},
	{Section: "submodule", Key: "branch", Type: "string", Default: "", DescZh: "子模块跟踪的分支", DescEn: "Branch tracked by the submodule"},

	// ========== tag ==========
	{Section: "tag", Key: "sort", Type: "string", Default: "-version:refname", DescZh: "标签排序方式（如 -version:refname 按版本号倒序）", DescEn: "Tag sorting order (e.g. -version:refname for version desc)"},

	// ========== init ==========
	{Section: "init", Key: "defaultBranch", Type: "string", Default: "master", DescZh: "新仓库的默认分支名（现代推荐 main）", DescEn: "Default branch name for new repos (modern: main)"},

	// ========== color ==========
	{Section: "color", Key: "ui", Type: "enum", Default: "auto", EnumValues: []string{"true", "false", "auto", "always"}, DescZh: "是否启用彩色输出：auto=终端支持时启用", DescEn: "Enable colored output: auto=enable when terminal supports"},
}

// GetKnownKeys 返回全部 KnownKey 定义
func GetKnownKeys() []KnownKey {
	return knownKeys
}

// GetKnownKeysForSection 返回指定节的 KnownKey 定义
func GetKnownKeysForSection(section string) []KnownKey {
	var result []KnownKey
	for _, k := range knownKeys {
		if k.Section == section {
			result = append(result, k)
		}
	}
	return result
}

// FindKnownKey 查找指定节和键名的 KnownKey
func FindKnownKey(section, key string) *KnownKey {
	for i, k := range knownKeys {
		if k.Section == section && k.Key == key {
			return &knownKeys[i]
		}
	}
	return nil
}
