// Package makefile
// @description: Makefile Visual Editor — 内置模板定义
package makefile

// GetBuiltinTemplates 返回所有内置 Makefile 模板。
// 内置模板涵盖：Go 多平台构建、Docker 镜像构建与推送、服务部署与重启、前端构建、通用清理。
func GetBuiltinTemplates() []Template {
	return []Template{
		{
			ID:          "builtin-go-multiplatform",
			Name:        "Go 多平台构建",
			Description: "编译 linux/windows/darwin 三平台二进制文件",
			IsBuiltin:   true,
			Doc: MakefileDoc{
				Variables: []Variable{
					{Name: "BINARY", Operator: ":=", Value: "app"},
					{Name: "VERSION", Operator: ":=", Value: "1.0.0"},
				},
				Targets: []Target{
					{
						Name:     "build",
						Deps:     []string{"build-linux", "build-windows", "build-darwin"},
						Commands: []string{},
						IsPhony:  true,
					},
					{
						Name:     "build-linux",
						Deps:     []string{},
						Commands: []string{"GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 ."},
						IsPhony:  true,
					},
					{
						Name:     "build-windows",
						Deps:     []string{},
						Commands: []string{"GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe ."},
						IsPhony:  true,
					},
					{
						Name:     "build-darwin",
						Deps:     []string{},
						Commands: []string{"GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 ."},
						IsPhony:  true,
					},
					{
						Name:     "clean",
						Deps:     []string{},
						Commands: []string{"rm -rf dist/"},
						IsPhony:  true,
					},
				},
				RawBlocks: []RawBlock{},
			},
		},
		{
			ID:          "builtin-docker",
			Name:        "Docker 镜像构建与推送",
			Description: "构建 Docker 镜像并推送到镜像仓库",
			IsBuiltin:   true,
			Doc: MakefileDoc{
				Variables: []Variable{
					{Name: "IMAGE", Operator: ":=", Value: "myapp"},
					{Name: "TAG", Operator: ":=", Value: "latest"},
					{Name: "REGISTRY", Operator: ":=", Value: "registry.example.com"},
				},
				Targets: []Target{
					{
						Name:     "docker-build",
						Deps:     []string{},
						Commands: []string{"docker build -t $(REGISTRY)/$(IMAGE):$(TAG) ."},
						IsPhony:  true,
					},
					{
						Name:     "docker-push",
						Deps:     []string{"docker-build"},
						Commands: []string{"docker push $(REGISTRY)/$(IMAGE):$(TAG)"},
						IsPhony:  true,
					},
					{
						Name:     "docker-clean",
						Deps:     []string{},
						Commands: []string{"docker rmi $(REGISTRY)/$(IMAGE):$(TAG)"},
						IsPhony:  true,
					},
				},
				RawBlocks: []RawBlock{},
			},
		},
		{
			ID:          "builtin-deploy",
			Name:        "服务部署与重启",
			Description: "将服务部署到远程服务器并重启",
			IsBuiltin:   true,
			Doc: MakefileDoc{
				Variables: []Variable{
					{Name: "SERVER", Operator: ":=", Value: "user@host"},
					{Name: "DEPLOY_DIR", Operator: ":=", Value: "/opt/app"},
					{Name: "SERVICE", Operator: ":=", Value: "myapp"},
				},
				Targets: []Target{
					{
						Name:     "deploy",
						Deps:     []string{"build"},
						Commands: []string{
							"scp dist/$(SERVICE) $(SERVER):$(DEPLOY_DIR)/",
							"ssh $(SERVER) 'systemctl restart $(SERVICE)'",
						},
						IsPhony: true,
					},
					{
						Name:     "restart",
						Deps:     []string{},
						Commands: []string{"ssh $(SERVER) 'systemctl restart $(SERVICE)'"},
						IsPhony:  true,
					},
					{
						Name:     "build",
						Deps:     []string{},
						Commands: []string{"go build -o dist/$(SERVICE) ."},
						IsPhony:  true,
					},
				},
				RawBlocks: []RawBlock{},
			},
		},
		{
			ID:          "builtin-frontend",
			Name:        "前端项目构建",
			Description: "使用 npm 或 yarn 构建前端项目",
			IsBuiltin:   true,
			Doc: MakefileDoc{
				Variables: []Variable{
					{Name: "PKG_MANAGER", Operator: ":=", Value: "npm"},
				},
				Targets: []Target{
					{
						Name:     "install",
						Deps:     []string{},
						Commands: []string{"$(PKG_MANAGER) install"},
						IsPhony:  true,
					},
					{
						Name:     "build",
						Deps:     []string{"install"},
						Commands: []string{"$(PKG_MANAGER) run build"},
						IsPhony:  true,
					},
					{
						Name:     "dev",
						Deps:     []string{"install"},
						Commands: []string{"$(PKG_MANAGER) run dev"},
						IsPhony:  true,
					},
					{
						Name:     "clean",
						Deps:     []string{},
						Commands: []string{"rm -rf dist/ node_modules/"},
						IsPhony:  true,
					},
				},
				RawBlocks: []RawBlock{},
			},
		},
		{
			ID:          "builtin-clean",
			Name:        "通用清理",
			Description: "清理构建产物和临时文件",
			IsBuiltin:   true,
			Doc: MakefileDoc{
				Variables: []Variable{},
				Targets: []Target{
					{
						Name:     "clean",
						Deps:     []string{},
						Commands: []string{"rm -rf dist/ build/ *.tmp"},
						IsPhony:  true,
					},
					{
						Name:     "clean-all",
						Deps:     []string{"clean"},
						Commands: []string{"rm -rf vendor/ node_modules/"},
						IsPhony:  true,
					},
				},
				RawBlocks: []RawBlock{},
			},
		},
	}
}
