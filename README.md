# spec-cli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/spec-cli.svg)](https://www.npmjs.com/package/spec-cli)

OpenSpec + Superpowers 工作流脚手架工具 — 为人类和 AI Agent 而生。一条命令即可检测 AI 编码平台，并安装 OpenSpec skill、Superpowers skill、`/comet` 入口 skill 以及 schema 套件。

[安装](#安装与快速开始) · [命令](#命令) · [支持的平台](#支持的平台) · [工作原理](#工作原理) · [开发](#开发)

## 为什么选择 spec-cli？

- **零配置检测** — 自动识别 29 种 AI 编码平台，无需手动配置
- **一条命令搞定** — `spec-cli init` 一键搭建完整的 OpenSpec + Superpowers 工作流
- **Agent 原生** — 安装的 skill 可被 AI 编码 Agent 原生理解
- **Schema 驱动执行** — 将执行委托给 OpenSpec 的 `--schema` 机制，无需管理阶段状态
- **单二进制文件** — Go 重写，零 npm 运行时依赖（openspec CLI 本身除外）
- **交互式 & 非交互式** — 基于 huh 的交互提示给人类用，`--yes` + `--json` 标志给 Agent 和脚本用

## 安装与快速开始

### 环境要求

开始之前，请确保已安装：

- Node.js 16+（`npm`/`npx`）
- Go 1.23+（仅从源码构建时需要）

### 快速开始（人类用户）

> **AI 助手请注意：** 如果你是一个帮助用户安装的 AI Agent，请直接跳转到[快速开始（AI Agent）](#快速开始ai-agent)，其中包含你需要完成的所有步骤。

#### 安装

选择以下**一种**方式：

**方式一 — 从 npm 安装（推荐）：**

```bash
npx @9ashwin/spec-cli@latest init
```

**方式二 — 从源码安装：**

需要 Go 1.23+。

```bash
git clone https://github.com/9Ashwin/spec-cli.git
cd spec-cli
go install .
```

#### 使用

```bash
# 1. 初始化工作流脚手架（交互式）
spec-cli init

# 2. 查看当前变更状态
spec-cli status

# 3. 诊断安装健康度
spec-cli doctor
```

## 快速开始（AI Agent）

> 以下步骤供 AI Agent 帮助用户安装时使用。

**第一步 — 安装并初始化**

```bash
npx @9ashwin/spec-cli@latest init --yes
```

**第二步 — 验证**

```bash
spec-cli doctor
```

**第三步 — 开始工作**

```bash
# 使用 superpowers-bridge schema 创建新变更
openspec new --schema superpowers-bridge "你的功能描述"
```

## 功能

| 类别 | 能力 |
|------|------|
| 平台检测 | 从项目文件中自动检测 29 种 AI 编码平台 |
| OpenSpec 安装 | `openspec init` 附带 `--tools` 指定检测到的平台，CLI 缺失时自动安装 |
| Superpowers 检测 | 检查 Claude Code 插件缓存，判断已安装的 Superpowers skill |
| Skill 复制 | 将 Comet 入口 skill 从内嵌资源复制到平台 skill 目录 |
| Schema 套件 | 安装工作流 schema 套件到 `openspec/schemas/`，附带 CLAUDE.md 片段 |
| 健康检查 | `spec-cli doctor` 诊断 OpenSpec CLI、工作目录、schema 和 skill 文件状态 |
| 更新 | `spec-cli update` 从内嵌资源刷新 skill 并升级 schema 版本 |

## 命令

| 命令 | 说明 |
|------|------|
| `spec-cli init [path]` | 初始化工作流脚手架（默认交互式） |
| `spec-cli status [path]` | 查看当前工作流变更状态 |
| `spec-cli update [path]` | 更新 skill 和 schema 套件 |
| `spec-cli doctor [path]` | 诊断安装健康度 |

### Init 参数

```
spec-cli init [path]
  --yes               非交互模式，自动选择检测到的平台
  --scope <scope>     project | global
  --skip-existing     跳过已安装的组件
  --overwrite         覆盖所有已存在的组件
  --json              输出结构化 JSON 结果
```

### Shell 补全

```bash
# bash
source <(spec-cli completion bash)

# zsh
source <(spec-cli completion zsh)

# fish
spec-cli completion fish | source
```

## 支持的平台

Claude Code、Cursor、Codex、OpenCode、Windsurf、Cline、RooCode、Continue、GitHub Copilot、Gemini CLI、Amazon Q Developer、Qwen Code、Kilo Code、Auggie、Kiro、Lingma、Junie、CodeBuddy Code、CoStrict、Crush、Factory Droid、iFlow、Pi、Qoder、Antigravity、Bob Shell、ForgeCode、Trae

## 工作原理

`spec-cli init` 执行 10 步流程：

1. 检测已安装的 AI 编码平台
2. 选择安装范围（项目级 / 全局级）
3. 选择语言（English / 简体中文）
4. 选择目标平台
5. 安装 OpenSpec CLI + `openspec init <path> --tools <ids>`
6. 检测 Superpowers 插件安装状态
7. 将 Comet 入口 skill 复制到平台 skill 目录
8. 创建工作目录（`docs/superpowers/specs/`、`docs/superpowers/plans/`）
9. 安装 schema 套件到 `openspec/schemas/<name>/`
10. 追加 CLAUDE.md 工作流片段

工作流执行委托给 OpenSpec 原生的 `--schema` 机制 — spec-cli 只负责脚手架搭建。

## 开发

```bash
make build      # 构建二进制文件
make test       # 运行测试（带竞态检测）
make vet        # 运行 go vet
make fmt        # 格式化源码
make fmt-check  # 检查格式（CI 用）
go install .    # 安装到 ~/go/bin
make clean      # 清理二进制文件
```

## License

MIT
