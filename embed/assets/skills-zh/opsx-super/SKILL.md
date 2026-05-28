---
name: opsx:super
description: "当用户要通过 /opsx:super 启动、路由或继续 OpenSpec + Superpowers 工作时使用。"
---

# opsx:super

把它作为 OpenSpec + Superpowers 工作的统一入口。它负责路由：判断是否需要 schema change，避免重复创建 active change，然后把执行交给 OpenSpec 的 `superpowers-bridge` schema。

<EXTREMELY-IMPORTANT>
这是一个 GATED 工作流。每个阶段（brainstorm, proposal, design, specs, tasks, plan, apply, verify）在 schema artifact 指令中定义了 exit gate，你必须遵守每一个 gate。每当 gate 需要用户输入时必须暂停——尤其是 brainstorm 审批、design 确认、范围变更、验证失败处理、分支/PR 决策。

第一个硬性 gate 是 brainstorm → proposal。brainstorm 的 EXIT GATE 写着："Stop here until the user has approved the proposed design direction." 在用户明确批准之前，不要写 proposal.md，不要推进下一步。

仅在阶段结果明确且 exit gate 已满足时才自动推进。schema exit gate 说停就停，停下来问用户。
</EXTREMELY-IMPORTANT>

## 指令优先级

用户指令始终优先。如果用户明确说不要创建 schema change，不要用本 skill 覆盖用户；说明取舍，然后按用户要求走。

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

创建或选择 change 后，进入下方的 Continuous Execution 循环。仅在 exit gate 满足时才自动推进——第一个 gate（brainstorm 审批）最关键，不要跳过它。

每次 `opsx:super` 调用进入 `superpowers-bridge` 后，都要先检查 OpenSpec status：

```bash
openspec status --change "<name>" --json
```

然后读取当前的 schema artifact 指令，从下一个未完成的 schema 步骤继续。推进到下一阶段前，验证当前阶段的 EXIT GATE（定义在 schema artifact 指令中）已满足。gate 需要用户输入就停止并询问。已完成的明确步骤自动推进：

brainstorm -> proposal -> design -> specs -> tasks -> plan -> apply action -> verify -> retrospective/archive。

每个 schema step 由 schema instructions 决定适用哪个 Superpowers skill。先调用相关 Superpowers skill，再执行该 schema step。不要在 schema step 要求 skill 时凭记忆手写 artifact。

Skill 输出路径覆盖（schema artifact 指令优先于 skill 默认值）：
- **brainstorming**：该 skill 的默认输出路径是 `docs/superpowers/specs/`，默认终端状态是调用 `writing-plans`。在 superpowers-bridge 下运行时，两者均需覆盖：将原始 brainstorming 输出写入 change 的 `brainstorm.md`（按 schema artifact 指令），用户审批设计方案后进入 **proposal** artifact——不要调用 writing-plans。
- **writing-plans**：该 skill 的默认输出路径是 `docs/superpowers/plans/`。在 superpowers-bridge 下运行时，将计划写入 change 的 `plan.md`。

artifact step 要从 OpenSpec 取得具体 instructions：

```bash
openspec instructions <artifact-id> --change "<name>" --json
```

使用返回的 `instruction`、`template`、`outputPath` 和 `dependencies`。把 artifact 写到 `outputPath`，然后重新运行 status 再推进。

apply action 不要在 status 里寻找 `apply` artifact。要从 OpenSpec 取得 action instructions：

```bash
openspec instructions apply --change "<name>" --json
```

使用返回的 `contextFiles`、`tasks`、`progress` 和 `instruction` 来执行实现阶段。

不要把连续执行委托给 `openspec-continue-change` 或 `/opsx:continue`；那个 skill 设计上会在一个 artifact 后停止。`opsx:super` 自己负责 status -> instructions -> artifact/action -> status 循环，直到遇到下面的停止条件。

把它当作“带门禁的连续工作流”，不是无条件一跑到底的 prompt：

- 只有当前 artifact/action 满足 schema instruction 和 exit gate 时，才自动推进。
- 恢复时以 OpenSpec status/instructions 和磁盘文件为准，不以对话历史为准。
- schema instruction 要求用户决策时必须停止并询问。不要替用户默认选择设计确认、范围扩大、验证失败处理或 branch/PR 处理。

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

创建后进入连续执行，推进这个新 change。

## 遵守 Bridge 路由

不要写入 `docs/superpowers/specs/`。
不要写入 `docs/superpowers/plans/`。
schema 产物应写入 `openspec/changes/<name>/`。

如果必需的 schema step 或 Superpowers skill 缺失，不要静默退回到裸 `brainstorming`、`writing-plans` 或手写 artifacts。停止并告诉用户缺了什么。

## 危险信号

这些想法表示你正在自我合理化。停止，重新读取 OpenSpec status 和 schema instructions 后再继续：

- “这是功能，但很小，可以跳过 schema。”
- “我先直接跑 brainstorming，之后再移动文件。”
- “已经有 active change，但新建一个更快。”
- “OpenSpec 失败了，所以我手动建目录。”
- “schema step 不清楚，所以我自己临时写下一个 artifact。”
