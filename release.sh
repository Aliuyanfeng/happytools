#!/bin/bash

# HappyTools 发布脚本
# 使用 git-chglog 自动生成 CHANGELOG 并发布

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印帮助信息
print_help() {
    echo "HappyTools 发布脚本"
    echo ""
    echo "用法: ./release.sh <版本号> [选项]"
    echo ""
    echo "示例:"
    echo "  ./release.sh v1.0.0          # 发布 v1.0.0"
    echo "  ./release.sh v1.0.0-beta.1   # 发布预发布版本"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -d, --dry-run  只预览，不执行实际操作"
    echo ""
}

# 检查参数
if [ $# -eq 0 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    print_help
    exit 0
fi

VERSION=$1
DRY_RUN=false

# 检查是否为 dry-run 模式
if [ "$2" = "-d" ] || [ "$2" = "--dry-run" ]; then
    DRY_RUN=true
fi

# 验证版本号格式
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo -e "${RED}❌ 错误: 版本号格式不正确${NC}"
    echo "期望格式: v1.0.0 或 v1.0.0-beta.1"
    exit 1
fi

echo -e "${GREEN}🚀 准备发布 $VERSION${NC}"
echo ""

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}❌ 错误: 有未提交的更改${NC}"
    echo "请先提交所有更改后再发布"
    exit 1
fi

# 检查标签是否已存在
if git tag | grep -q "^${VERSION}$"; then
    echo -e "${RED}❌ 错误: 标签 $VERSION 已存在${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 步骤 1/4: 生成 CHANGELOG（使用 --next-tag）${NC}"
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 将生成 CHANGELOG.md（预览 $VERSION）"
else
    # 使用 --next-tag 预览新标签的内容
    git-chglog --next-tag "$VERSION" -o CHANGELOG.md
    echo -e "${GREEN}✅ CHANGELOG.md 已生成（包含 $VERSION 的变更）${NC}"
fi

echo ""
echo -e "${YELLOW}📝 步骤 2/4: 提交 CHANGELOG${NC}"
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 将提交 CHANGELOG.md"
else
    git add CHANGELOG.md
    git commit -m "docs: update CHANGELOG for $VERSION"
    echo -e "${GREEN}✅ CHANGELOG 已提交${NC}"
fi

echo ""
echo -e "${YELLOW}🏷️  步骤 3/4: 创建标签${NC}"
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 将创建标签 $VERSION"
else
    git tag "$VERSION"
    echo -e "${GREEN}✅ 标签 $VERSION 已创建${NC}"
fi

echo ""
echo -e "${YELLOW}📤 步骤 4/4: 推送代码和标签${NC}"
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 将推送代码和标签到远程仓库"
else
    # 推送代码
    git push origin main
    # 推送标签
    git push origin "$VERSION"
    echo -e "${GREEN}✅ 代码和标签已推送${NC}"
fi

echo ""
if [ "$DRY_RUN" = true ]; then
    echo -e "${YELLOW}🎉 预览完成！${NC}"
    echo "移除 -d 或 --dry-run 参数以执行实际发布"
else
    echo -e "${GREEN}🎉 发布成功！${NC}"
    echo ""
    echo "📦 GitHub Actions 将自动构建和发布应用"
    echo "🔗 查看进度: https://github.com/Aliuyanfeng/happytools/actions"
    echo ""
    echo "📝 发布说明将自动从 CHANGELOG.md 提取"
fi
