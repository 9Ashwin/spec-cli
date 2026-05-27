---
name: comet
description: "Comet — OpenSpec + Superpowers 开发工作流。"
---

# Comet — OpenSpec + Superpowers 开发工作流

## 快速开始

```bash
# 创建新变更
openspec new --schema superpowers-bridge "你的想法"

# 查看活跃变更
openspec list --json

# 完成后归档
openspec archive --change "<name>" -y
```

## 工作流

OpenSpec 负责 WHAT（提案、spec 生命周期、归档）。
Superpowers 负责 HOW（头脑风暴、设计、计划、构建、验证）。

`superpowers-bridge` schema 将两者桥接起来 —— 按每步注入的 artifact 指令操作即可。
