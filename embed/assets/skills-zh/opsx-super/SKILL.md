---
name: opsx:super
description: "当用户要通过 /opsx:super 启动、路由或继续 OpenSpec + Superpowers 工作时使用。"
---

# opsx:super

把它作为 OpenSpec + Superpowers 工作的统一入口。它负责路由：判断是否需要 schema change，避免重复创建 active change，然后把执行交给 OpenSpec 的 `superpowers-bridge` schema。

## 指令优先级

路由工作时按以下顺序处理：

1. 用户明确指令
2. 项目指令，例如 AGENTS.md / CLAUDE.md
3. OpenSpec schema artifact instructions
4. 本入口 skill
5. 默认 agent 行为

如果这些指令冲突，停止并说明冲突，不要静默选择其中一条路径。

## 核心规则

创建 change 之前，必须先做两件事：

1. 判断请求属于 schema change 还是 direct PR。
2. 用下面命令检查 active changes：

```bash
openspec list --json
```

不要静默创建重复 change。

如果 `openspec` 缺失、项目未初始化，或 `superpowers-bridge` 未安装，停止并告诉用户运行 `spec-cli init` 或 `spec-cli doctor`。

## 路由请求

以下情况使用 `superpowers-bridge`：

- 新功能或新能力
- 架构变更
- 破坏性变更
- 对外契约、schema、数据模型或跨系统变更

以下情况不要创建 schema change：

- 不改变契约、只恢复预期行为的 bug fix
- typo 或纯文档修改
- 测试补齐
- linter / config value 微调
- 非破坏性依赖升级

如果属于 direct PR 场景，告诉用户这不需要 `opsx:super` change，然后按普通开发流程处理。

## Active Change 处理

- 没有 active change：如果请求需要 change，就创建新 change。
- 只有一个 active change：询问用户是继续它还是创建新 change，除非用户明确要求创建新 change。
- 有多个 active changes：列出它们，询问继续哪一个，或是否创建新 change。
- 用户要求继续：按 OpenSpec status / instructions 推进已有 change，不要创建新目录。

## 连续执行

不要在创建或选择 change 后停下。

每次 `opsx:super` 调用进入 `superpowers-bridge` 后，都要检查 OpenSpec status 和 schema artifact instructions，然后从下一个未完成的 schema step 继续。对无歧义的步骤自动推进：

brainstorm -> proposal -> design -> specs -> tasks -> plan -> apply -> verify -> retrospective/archive。

只有遇到以下情况才停止：

- schema workflow 已完成。
- 下一个 schema instruction 需要用户明确决策。
- OpenSpec status/instructions 缺失，或与项目/用户指令冲突。
- 必需的 schema step、Superpowers skill、命令或 artifact 缺失。
- 验证失败，或 schema instruction 要求停止。

不要把对话历史当作进度的事实来源。恢复上下文、工具失败后、推进到下一阶段前，都要重新读取 OpenSpec status/instructions。

## 创建新 Change

将用户需求转成简短 kebab-case 名称，然后执行：

```bash
openspec new change "<kebab-case-name>" --schema superpowers-bridge --description "<原始需求>"
```

创建后按 schema artifact instructions 推进：

brainstorm -> proposal -> design -> specs -> tasks -> plan -> apply -> verify -> retrospective/archive。

## 遵守 Bridge 路由

不要写入 `docs/superpowers/specs/`。
不要写入 `docs/superpowers/plans/`。
schema 产物应写入 `openspec/changes/<name>/`。

如果必需的 schema step 或 Superpowers skill 缺失，不要静默退回到裸 `brainstorming`、`writing-plans` 或手写 artifacts。停止并告诉用户缺了什么。

## 危险信号

发现以下想法时，停止并重新路由：

- “这是功能，但很小，可以跳过 schema。”
- “我先直接跑 brainstorming，之后再移动文件。”
- “已经有 active change，但新建一个更快。”
- “OpenSpec 失败了，所以我手动建目录。”
- “schema step 不清楚，所以我自己临时写下一个 artifact。”
