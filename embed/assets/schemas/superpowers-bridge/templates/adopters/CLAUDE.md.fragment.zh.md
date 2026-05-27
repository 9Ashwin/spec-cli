<!-- Source: superpowers-bridge/templates/adopters/CLAUDE.md.fragment.zh.md -->
<!-- 把这一节粘贴进你项目的 CLAUDE.md,让 Claude 知道如何分流本 repo 使用本 schema 的工作。 -->
<!-- 若你自定义 schema 名称或 bridge repo URL,请对应修改;否则保持原样即可。 -->

## 变更工作流(Claude Code 启动先读)

本 repo 采用 [`superpowers-bridge`](https://github.com/JiangWay/openspec-schemas/tree/main/superpowers-bridge) 对接 OpenSpec 与 Superpowers。整合规则(语言、artifact 路径、PRECHECK)以该 bridge README 为准;以下是给 Claude 的 routing 指引。

### 入口分流

| 你看到的触发 | 应该如何做 |
|---|---|
| 用户以 narrative 开“设计讨论 / 头脑风暴” | 先 verbal `superpowers:brainstorming`,**不**写到 `docs/superpowers/specs/`;对话收敛后依下方 5 条判断标准升级到 `/opsx:propose` |
| 用户直接调用 `/opsx:new` / `/opsx:ff` / `/opsx:propose` | 按 schema 规定的流程;artifact instruction 会在每步注入 |
| 用户明确说 bug fix / typo / config 微调 / 文件更新 | 直接 PR,**不**建 change(见下方 skip 规则) |
| 已经在某个 change 中 | `/opsx:continue` 或 `/opsx:apply` / `/opsx:verify` / `/opsx:archive` 持续推进 |

### 连续但带门禁的执行

进入 `superpowers-bridge` change 后,以 OpenSpec status 和 artifact
instructions 作为进度事实来源。已经满足输入和 exit gate 的阶段可以自动推进,
但遇到这些明确决策点必须暂停:设计方向未确认、缺少必要 skills/tools、
capability 范围不清楚、design 有阻塞问题、plan 扩大范围、verification
失败、spec/design drift 改变实现范围、最终 branch/PR 处理。

恢复时,重新运行 `openspec status --change <name> --json`,并检查现有
artifacts。不要根据聊天历史推断进度。

### 何时**不**按 opsx(直接 PR)

| 情境 | 直接 PR? |
|---|---|
| 新功能 / 新 capability / 架构变更 / breaking change | ❌ 要按 opsx |
| Bug fix(不变更合约)/ 测试补写 / linter 规则 / 非破坏性升级 / typo / 文件 / config 值微调 | ✅ 直接 PR |

原则:**流程仪式跟风险成正比**。动到对外合约 / schema / 跨系统对接 / 合规边界 → opsx;其他 → 直接 PR。

### Verbal brainstorm 升级到 opsx 的 5 条判断标准

5 条**全部满足**才升级(任一不满足则继续 brainstorm,不写到 `docs/superpowers/specs/`):

1. **Scope 锁定** —— 一句话讲清“包含/不包含什么”
2. **主要设计分歧已收敛** —— 替代方案选过,剩下 TBD 有明确 owner 与影响面
3. **跨系统依赖盘点过** —— 对方就绪 / 暂 mock / 真未知,三选一说明清楚
4. **验收条件可陈述** —— 具体 pass 条件(例:`./mvnw clean verify` 通过 + N 个成果)
5. **对话进入收敛** —— 最近几轮在 confirm 不在发散

全部满足 → 主动建议用户“要不要 `/opsx:propose`?”,用户确认后落地。永远不要自动触发。

### Front-door 反模式(别做)

- 让 brainstorming 写到 `docs/superpowers/specs/`
- 让 writing-plans 写到 `docs/superpowers/plans/`
- TBD 没收敛就升级到 opsx
- 对 bug fix / typo 也建 change

详细见 [superpowers-bridge README §进入与离开的判断](https://github.com/JiangWay/openspec-schemas/blob/main/superpowers-bridge/README.zh.md#进入与离开的判断entry--exit-gates)。
