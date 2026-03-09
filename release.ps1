# HappyTools 发布脚本 (Windows PowerShell)
# 使用 git-chglog 自动生成 CHANGELOG 并发布

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$Version,

    [Parameter(Mandatory=$false)]
    [switch]$DryRun
)

# 颜色输出函数
function Write-Success { Write-Host "✅ $args" -ForegroundColor Green }
function Write-Error { Write-Host "❌ $args" -ForegroundColor Red }
function Write-Info { Write-Host "ℹ️  $args" -ForegroundColor Cyan }
function Write-Step { Write-Host "📋 $args" -ForegroundColor Yellow }

# 打印帮助信息
if ($Version -eq "-h" -or $Version -eq "--help") {
    Write-Host "HappyTools 发布脚本"
    Write-Host ""
    Write-Host "用法: .\release.ps1 <版本号> [-DryRun]"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  .\release.ps1 v1.0.0          # 发布 v1.0.0"
    Write-Host "  .\release.ps1 v1.0.0-beta.1   # 发布预发布版本"
    Write-Host "  .\release.ps1 v1.0.0 -DryRun  # 预览发布流程"
    Write-Host ""
    exit 0
}

# 验证版本号格式
if ($Version -notmatch '^v\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$') {
    Write-Error "版本号格式不正确"
    Write-Host "期望格式: v1.0.0 或 v1.0.0-beta.1"
    exit 1
}

Write-Host "🚀 准备发布 $Version" -ForegroundColor Green
Write-Host ""

# 检查是否有未提交的更改
$gitStatus = git status --porcelain
if ($gitStatus) {
    Write-Error "有未提交的更改"
    Write-Host "请先提交所有更改后再发布"
    exit 1
}

# 检查标签是否已存在
$tagExists = git tag | Select-String "^$Version$"
if ($tagExists) {
    Write-Error "标签 $Version 已存在"
    exit 1
}

Write-Step "步骤 1/4: 生成 CHANGELOG（使用 --next-tag）"
if ($DryRun) {
    Write-Info "[DRY-RUN] 将生成 CHANGELOG.md（预览 $Version）"
} else {
    # 使用 --next-tag 预览新标签的内容
    git-chglog --next-tag $Version -o CHANGELOG.md
    Write-Success "CHANGELOG.md 已生成（包含 $Version 的变更）"
}

Write-Host ""
Write-Step "步骤 2/4: 提交 CHANGELOG"
if ($DryRun) {
    Write-Info "[DRY-RUN] 将提交 CHANGELOG.md"
} else {
    git add CHANGELOG.md
    git commit -m "docs: update CHANGELOG for $Version"
    Write-Success "CHANGELOG 已提交"
}

Write-Host ""
Write-Step "步骤 3/4: 创建标签"
if ($DryRun) {
    Write-Info "[DRY-RUN] 将创建标签 $Version"
} else {
    git tag $Version
    Write-Success "标签 $Version 已创建"
}

Write-Host ""
Write-Step "步骤 4/4: 推送代码和标签"
if ($DryRun) {
    Write-Info "[DRY-RUN] 将推送代码和标签到远程仓库"
} else {
    # 推送代码
    git push origin main
    # 推送标签
    git push origin $Version
    Write-Success "代码和标签已推送"
}

Write-Host ""
if ($DryRun) {
    Write-Host "🎉 预览完成！" -ForegroundColor Yellow
    Write-Host "移除 -DryRun 参数以执行实际发布"
} else {
    Write-Host "🎉 发布成功！" -ForegroundColor Green
    Write-Host ""
    Write-Host "📦 GitHub Actions 将自动构建和发布应用"
    Write-Host "🔗 查看进度: https://github.com/Aliuyanfeng/happytools/actions"
    Write-Host ""
    Write-Host "📝 发布说明将自动从 CHANGELOG.md 提取"
}
