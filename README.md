# spec-cli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25.7-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/ashwin-spec.svg)](https://www.npmjs.com/package/ashwin-spec)

OpenSpec + Superpowers 工作流脚手架工具 — 为人类和 AI Agent 而生。一条命令即可检测 AI 编码平台，并安装 OpenSpec skill、`/opsx:super` 入口 skill 以及 schema 资源包。

[安装](#安装与快速开始) · [命令](#命令) · [支持的平台](#支持的平台) · [工作原理](#工作原理) · [开发](#开发)

## 为什么选择 spec-cli？

- **零配置检测** — 自动识别 29 种 AI 编码平台，无需手动配置
- **一条命令搭好脚手架** — `spec-cli init` 一键安装 OpenSpec skill、入口 skill 和 schema 资源包
- **Agent 原生** — 安装的 skill 可被 AI 编码 Agent 原生理解
- **Schema 驱动执行** — 将执行委托给 OpenSpec 的 `--schema` 机制，无需管理阶段状态
- **单二进制文件** — Go CLI，npm 包只负责下载和启动平台对应的二进制文件（openspec CLI 本身除外）
- **交互式 & 非交互式** — 基于 huh 的交互提示给人类用，`--yes` + `--json` 标志给 Agent 和脚本用

## 安装与快速开始

### 环境要求

开始之前，请确保已安装：

- Node.js 16+（`npm`/`npx`）
- Go 1.25.7+（仅从源码构建时需要）

### 快速开始（人类用户）

> **AI 助手请注意：** 如果你是一个帮助用户安装的 AI Agent，请直接跳转到[快速开始（AI Agent）](#快速开始ai-agent)，其中包含你需要完成的所有步骤。

#### 安装

选择以下**一种**方式：

**方式一 — 全局安装（推荐）：**

```bash
npm install -g ashwin-spec@latest
```

**方式二 — 按需运行：**

```bash
npx ashwin-spec@latest init
```

**方式三 — 从源码安装：**

需要 Go 1.25.7+。

```bash
git clone https://github.com/9Ashwin/spec-cli.git
cd spec-cli
go install .
```

#### 使用

```bash
# 第一步：初始化工作流脚手架（交互式）
spec-cli init

# 第二步：在 AI 编码平台中输入 /opsx:super 开始新变更
# 也可以在终端直接创建：
openspec new change "your-change-name" --schema superpowers-bridge --description "你的需求描述"

# 第三步：让 /opsx:super 按 schema instructions 连续推进
# brainstorming → proposal → design → specs → tasks → plan → apply → verify → archive

# 查看进度
spec-cli status

# 健康检查
spec-cli doctor
```

### 完整使用流程

```
1. spec-cli init          — 搭建脚手架（一次性，安装 OpenSpec + skill + schema）

2. /opsx:super            — 在 AI 编码平台中启动新变更
   或 openspec new change — 终端直接创建

3. /opsx:super 连续执行 OpenSpec instructions
   brainstorming          — AI 跟你讨论需求、收敛设计方向
   proposal/design/specs  — 写出 schema artifacts
   tasks/plan             — 准备实现工作
   apply                  — 按 action instructions 实现并测试
   verify                 — 验证实现正确性
   archive                — 从当前 worktree 归档并同步 spec

4. spec-cli status         — 随时查看活跃变更
   spec-cli doctor         — 诊断环境是否正常
   spec-cli update         — 升级 skill 和 schema 版本
```

## 快速开始（AI Agent）

> 以下步骤供 AI Agent 帮助用户安装时使用。

**第一步 — 安装**

```bash
npm install -g ashwin-spec@latest
```

**第二步 — 初始化**

```bash
spec-cli init --yes
```

**第三步 — 验证**

```bash
spec-cli doctor
```

**第四步 — 开始工作**

```bash
# 使用 superpowers-bridge schema 创建新变更
openspec new change "your-change-name" --schema superpowers-bridge --description "你的功能描述"
```

## 功能

| 类别 | 能力 |
|------|------|
| 平台检测 | 从项目文件中自动检测 29 种 AI 编码平台 |
| OpenSpec 安装 | `openspec init` 附带 `--tools` 指定检测到的平台，CLI 缺失时尝试通过 npm 安装 |
| Superpowers 检测 | 检查 Claude Code 插件缓存，判断是否已安装 Superpowers skill |
| Skill 复制 | 将 opsx:super 入口 skill 从内嵌资源复制到平台 skill 目录 |
| Schema 资源包 | 安装工作流 schema 资源包到 `openspec/schemas/`，附带 CLAUDE.md 片段 |
| 健康检查 | `spec-cli doctor` 诊断 OpenSpec CLI、项目路径、schema 和 skill 文件状态 |
| 更新 | `spec-cli update` 从内嵌资源刷新 skill 并升级 schema 版本 |

## 命令

| 命令 | 说明 |
|------|------|
| `spec-cli init [path]` | 初始化工作流脚手架（默认交互式） |
| `spec-cli status [path]` | 查看当前工作流变更状态 |
| `spec-cli update [path]` | 更新 skill 和 schema 资源包 |
| `spec-cli doctor [path]` | 诊断安装健康度 |
| `spec-cli completion <shell>` | 生成 shell 补全脚本（bash/zsh/fish/powershell） |
| `spec-cli --version` | 查看版本（包含 commit hash 和 构建日期） |

### Init 参数

```
spec-cli init [path]
  --yes               非交互模式，自动选择检测到的平台
  --scope <scope>     project | global
  --skip-existing     跳过已安装的组件
  --overwrite         覆盖所有已存在的组件
  --json              输出结构化 JSON 结果
```

### Update 参数

```
spec-cli update [path]
  --scope <scope>       project | global
  --language <lang>     en | zh
  --json                输出结构化 JSON 结果
```

### Doctor 参数

```
spec-cli doctor [path]
  --scope <scope>       auto | project | global
  --json                输出结构化 JSON 结果
```

### Shell 补全

```bash
# bash
eval "$(spec-cli completion bash)"

# zsh
eval "$(spec-cli completion zsh)"

# fish
spec-cli completion fish | source

# powershell
spec-cli completion powershell | Out-String | Invoke-Expression
```

## 支持的平台

Claude Code、Cursor、Codex、OpenCode、Windsurf、Cline、RooCode、Continue、GitHub Copilot、Gemini CLI、Amazon Q Developer、Qwen Code、Kilo Code、Auggie、Kiro、Lingma、Junie、CodeBuddy Code、CoStrict、Crush、Factory Droid、iFlow、Pi、Qoder、Antigravity、Bob Shell、ForgeCode、Trae

## 工作原理

`spec-cli init` 执行 9 步流程：

1. 检测已安装的 AI 编码平台
2. 选择安装范围（项目级 / 全局级）
3. 选择语言（English / 简体中文）
4. 选择目标平台
5. 安装 OpenSpec CLI + `openspec init <path> --tools <ids>`
6. 检测 Superpowers 插件安装状态
7. 将 opsx:super 入口 skill 复制到平台 skill 目录
8. 安装 schema 资源包到 `openspec/schemas/<name>/`
9. 追加 CLAUDE.md 工作流片段

工作流执行委托给 OpenSpec 原生的 `--schema` 机制 — spec-cli 只负责脚手架搭建。

安装后的 `/opsx:super` skill 会保持轻量：读取 `openspec status --change "<name>" --json`，用 `openspec instructions ... --json` 获取下一个 artifact 或 action instructions，调用对应的 Superpowers skill，把 artifact 写到 OpenSpec 返回的 `outputPath`，然后重新检查 status 再推进。

准备归档时，请在包含最新已勾选 `tasks.md`、`verify.md`、`retrospective.md` 和 implementation commits 的同一个 branch/worktree 中运行 `openspec archive <change-name> -y`。不要从 stale checkout 归档，否则可能把旧的 task 状态同步进 spec。

## 开发

### 常用命令

```bash
make build      # 构建二进制文件（含 Version + Date + CommitHash 注入）
make test       # 运行测试（带竞态检测）
make vet        # 运行 go vet
make lint       # 运行 golangci-lint
make fmt        # 格式化源码
make fmt-check  # 检查格式（CI 用）
make release    # 构建 npm 分发归档和 checksums
go install .    # 安装到 ~/go/bin
make clean      # 清理二进制文件
```

### 开发环境设施

- **golangci-lint**：使用 `forbidigo` 强制 `internal/` 下通过 `vfs` 包中转所有文件系统调用
- **husky + lint-staged**：`pnpm install` 后自动启用 pre-commit hook（fmt-check + vet + gofmt）
- **gitleaks**：`.gitleaks.toml` 配置 secret 泄露检测，提交前可手动运行 `gitleaks detect`
- **GoReleaser**：`.goreleaser.yml` 驱动多平台发布（darwin/linux/windows × amd64/arm64）

### 架构亮点

- **VFS 抽象**：`internal/vfs/FS` 接口 — 所有文件操作流过可替换的实现，测试可用内存 FS
- **Runner 接口**：`internal/openspec/Runner` 抽象了 `os/exec`，`mockRunner` 使 openspec 集成全面可测
- **嵌入 FS 跨平台**：`embed.FS` 使用 `path.Join`（正斜杠），写入本地时转为 `filepath.Join`，避免 Windows 路径问题

## License

MIT
