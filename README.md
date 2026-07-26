# Copy-Paste AI 🤖📋

> **我们不生产 token，我们只是 token 的搬运工** (*^_^*)

[![CI](https://github.com/yuan-shuo/copy-paste-ai/workflows/ci/badge.svg)](https://github.com/yuan-shuo/copy-paste-ai/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuan-shuo/copy-paste-ai)](https://goreportcard.com/report/github.com/yuan-shuo/copy-paste-ai)
[![codecov](https://codecov.io/gh/yuan-shuo/copy-paste-ai/branch/main/graph/badge.svg)](https://codecov.io/gh/yuan-shuo/copy-paste-ai)
[![Release](https://img.shields.io/github/release/yuan-shuo/copy-paste-ai.svg)](https://github.com/yuan-shuo/copy-paste-ai/releases/latest)

---

以下内容大部分由AI帮我生成，因为我懒，如果你看完绷住了，你可以给我10万美元。

## 宇宙级项目介绍

**Copy-Paste AI（简称 `cpa`）** 是一款跨时代的革命性工具，具备以下宇宙级特性：

- 🌌 **兼容可观测宇宙中全部网页版 AI 大语言模型**（已在 ChatGPT、Kimi、元宝、豆包、文心一言、通义、Cloud 验证通过；理论上覆盖了银河系内所有碳基生物可能使用的 AI 服务）
- 🌠 **支持银河系任何一门编程语言**（包括但不限于：Go、Rust、Python、JavaScript、Java、C++、Zig、R、Haskell、Erlang、COBOL、Brainfuck、Malbolge、LOLCODE 以及木卫二冰层下可能存在的海底文明编程语言）
- 🌀 **全程无任何资费消耗的永动机**（不调 API、不走网关、不消耗 token、不产生碳排放。唯一的成本是你下载这个二进制的那几秒钟电费）
- 🧠 **零学习成本**（如果你会复制粘贴，你就会用这个工具。如果你不会复制粘贴，那你现在有两个问题需要解决）

**核心理念：** 让 AI 帮你写代码，但让你不用花钱、不用配 key、不用学 prompt engineering，只需拖一个文件到浏览器对话框里，然后躺平等结果。

---

## 安装

### 方式一：Go 安装（推荐）

```bash
go install github.com/yuan-shuo/copy-paste-ai@latest
```

### 方式二：下载二进制

去 [Releases](https://github.com/yuan-shuo/copy-paste-ai/releases) 下载对应平台的压缩包然后塞到你的环境变量目录下。目前支持：

- 🪟 Windows（amd64 / arm64）
- 🐧 Linux（amd64 / arm64）
- 🍎 macOS（amd64 / arm64）

> ⚠️ 其他平台（如 FreeBSD、OpenBSD、Plan 9）暂未支持，我们正在非常不积极地联系这些平台的使用者，请不要抱有任何幻想。

### 方式三：源码编译

```bash
git clone https://github.com/yuan-shuo/copy-paste-ai.git
cd copy-paste-ai
go build -o cpa .
```

---

## 快速上手：三步召唤 AI 之术

### 第一步：窥探项目全貌

```bash
cpa tree
```

```
my-awesome-project/
├── main.go
├── go.mod
├── cmd/
│   ├── gen.go
│   └── tree.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── gitignore/
│   │   └── gitignore.go
│   └── tree/
│       └── tree.go
└── assets/
    └── skel/
        ├── config.toml
        ├── prompt.tmpl
        └── prompt_default.md
```

> ✨ 自动忽略 `.git` 和 `.cpa` 目录，自动遵循 `.gitignore` 规则。你的项目就像被 AI 用望远镜拍了一张结构照。

### 第二步：打包召唤物资

```bash
# 指定你想让 AI 看的文件
cpa gen -f main.go,internal/config/config.go

# 或者用配置好的别名
cpa gen -f core

# 什么都不指定？也行，就给 AI 看个文件树
cpa gen
```

这会在 `.cpa/prompt/` 目录下生成一个 `.md` 文件，里面包含：

1. 📂 **项目文件树** — AI 看了会说"哦，原来你这个项目是这样的"
2. 📝 **你指定的文件代码** — AI 看了会说"嗯，这个代码写得...还行吧"
3. 💬 **你的定制提示词** — AI 看了会说"好的老板，我这就干活"

### 第三步：投弹行动 🎯

```
┌──────────────────────────────────────────────────────┐
│  🖱️  打开网页版 AI（ChatGPT / Kimi / 元宝 / ...）      │
│  📎  点击对话框的 📎 附件按钮，或者直接拖拽 🤙        │
│  📄  把 .cpa/prompt/20260726130000.md 拖进去        │
│  ✍️  输入你的需求："帮我重构这个函数"                │
│  🚀  按下回车，双手离开键盘                          │
│  🤖  看着 AI 开始工作，你可以去喝杯水 ☕              │
│  📥  把 AI 产出的代码粘回你的项目                     │
│  ✅  go build && go test → 验证通过                  │
│  🎉  收工！摸鱼时间                                  │
└──────────────────────────────────────────────────────┘
```

> **是的，现代 AI 网页对话框都支持直接上传/拖拽 md 文件了。** 我们活在一个了不起的时代。

---

## 命令详解

### `cpa tree` — 项目全景扫描

```bash
cpa tree
```

**特性：**

- 🚫 自动跳过 `.git/` 和 `.cpa/` 目录（这俩目录 AI 不需要看）
- 📋 默认遵循 `.gitignore` 规则（你本来就不想被 AI 看到的文件，现在真的看不到了）
- 📺 直接打印到终端，不生成任何文件，不写任何磁盘，环保 🌱

### `cpa gen` — AI 物资打包器

```bash
# 仅生成文件树
cpa gen

# 指定多个文件
cpa gen -f main.go,go.mod,internal/app/app.go

# 使用别名
cpa gen -f core

# 混合使用
cpa gen -f main,utils,extra
```

**参数：**

| 参数            | 说明                                     |
| --------------- | ---------------------------------------- |
| `-f, --files` | 要包含的文件列表，逗号分隔，支持文件别名 |

### `cpa --help` — 查看帮助

```bash
cpa --help
```

---

## 配置文件

第一次运行 `cpa gen` 时，会在项目根目录自动创建 `.cpa/config.toml`：

```toml
# ============================================================
# Copy-Paste AI 配置文件
# 所有配置默认都被注释掉了，取消注释即可启用
# 就像你人生中的那些默认选项一样，改不改由你
# ============================================================


# ============================================================
# [default] 默认文件
# 每次 cpa gen 都会自动包含这些文件，不用每次 -f 指定
# 用法: files = ["路径1", "路径2"]
# ============================================================

# [default]
# files = ["main.go", "go.mod"]


# ============================================================
# [file_aliases] 文件别名
# 给一组经常一起出现的文件起个短名字
# 用法: "别名" = ["文件1", "文件2"]
# 例子: cpa gen -f core 等价于 cpa gen -f pkg/core/core.go,pkg/core/utils.go
# ============================================================

# [file_aliases]
# "main" = ["main.go"]
# "test" = ["internal/test.go"]
# "core" = ["pkg/core/core.go", "pkg/core/utils.go"]


# ============================================================
# [prompt] AI 提示词
# 这就是你给 AI 下达的圣旨，想怎么写就怎么写
# 会被追加到生成的 md 文件末尾
# ============================================================

# [prompt]
# content = """
# # 指令
# 请根据以上代码完成我的需求。
# 要求：保持原有注释，总结所有改动。
# """


# ============================================================
# [gitignore] .gitignore 规则
# 是否尊重你项目里的 .gitignore 文件
# 默认: true（尊重）
# ============================================================

# [gitignore]
# enabled = true
```

### 配置项一览

| 配置项                | 类型                    | 说明                                   |
| --------------------- | ----------------------- | -------------------------------------- |
| `default.files`     | `[]string`            | 每次 gen 自动包含的文件，不用每次 -f   |
| `file_aliases`      | `map[string][]string` | 文件别名，`cpa gen -f core` 一键展开 |
| `prompt.content`    | `string`              | AI 提示词，追加到 md 末尾              |
| `gitignore.enabled` | `bool`                | 是否遵循 .gitignore，默认`true`      |

---

## 实战场景

### 场景一：让 AI 帮你重构一坨屎山代码

```bash
# 1. 先看看屎山的地形
cpa tree

# 2. 打包屎山核心文件
cpa gen -f internal/tree/tree.go

# 3. 拖文件给 AI，说："这个函数圈复杂度 31，帮我拆成 10 以下"
# 4. AI 开始工作，你去泡杯茶
# 5. AI 产出代码，粘回项目
# 6. go test ./... → 全绿！
# 7. 你已经成为了更好的工程师。或者说，AI 让你成为了更好的工程师。
```

### 场景二：让 AI 帮你写单元测试（AI 的传统手艺）

```bash
# 1. 打包待测文件
cpa gen -f internal/config/config.go

# 2. 拖给 AI，说："给这个文件写满单元测试，覆盖率 100%"
# 3. AI 噼里啪啦输出 200 行测试代码
# 4. 粘回去，go test -cover → 98%
# 5. 完美。那 2% 未覆盖的是 AI 故意留给你的，保持手感用的
```

### 场景三：让 AI 帮你读代码（当接手新项目时）

```bash
# 1. 生成全量上下文
cpa gen -f main.go,go.mod,cmd/gen.go,cmd/tree.go,internal/config/config.go

# 2. 拖给 AI，说："用通俗的语言给我讲一下这个项目是干嘛的"
# 3. AI 开始口若悬河
# 4. 你："哦原来如此！"
# 5. 3 分钟后，你已经比看 3 小时 README 理解得更透彻
```

---

## 宇宙级 FAQ

### Q: 为什么不直接调 AI API？

A: **因为我们是一个有态度的项目。** 真正的原因：

- 🔒 企业代理环境用不了 API
- 💰 个人用户不想花钱
- 🎫 不想到处配 key，怕泄露
- 🤷 有些 AI 服务根本没 API（比如某些网页版是独占的）
- 😇 网页版拖拽文件的体验已经足够好了，没必要重复造轮子

### Q: 生成的 md 文件 AI 能直接读吗？

A: **能。** 现在的 AI 网页对话框基本都支持拖拽/上传文件。md 文件格式简单，渲染器都认识。你甚至可以把多个 md 一起拖进去，AI 会同时处理。

### Q: 支持哪些语言的语法高亮？

A: **MD 文件用 GFM（GitHub Flavored Markdown），理论上支持 300+ 种语法高亮。** 实际上 AI 会把 md 里的 ``python / ``go / ```rust 等代码块当作对应语言来理解。哪怕你用 Brainfuck，AI 也认得。（也许吧）

### Q: 生成的文件会泄露项目信息吗？

A: **取决于你让 AI 看什么。** `cpa gen -f` 指定哪些文件就打包哪些文件。`.env`、`config.yaml` 这类敏感文件别加进去，这是常识。**工具不替你做安全决策，你得自己动脑子。**

### Q: 我删了 .cpa 目录会怎样？

A: **没事。** 下次 `cpa gen` 会自动重建。配置会重置。你可以放心删，大胆试。

### Q: `cpa tree` 和 `tree` 命令有啥区别？

A: 两个主要区别：

1. `tree` 是操作系统自带的，`cpa tree` 是我们写的
2. `tree` 不认识 `.gitignore`，会把 `node_modules/` 全给你打出来，显示 3 层你就想死了。`cpa tree` 完美遵循 `.gitignore`，干净清爽。

### Q: 能在 Windows 上用吗？

A: **当然能。** 我们就是在 Windows 上开发的。`go build` 一行搞定。

### Q: 这个项目有 bug 吗？

A: **肯定有。** 但我们相信 AI 能帮我们修。所以我们已经把自己逼入了"用 AI 修自己"的闭环中 🌀

---

## 项目结构

```
copy-paste-ai/
├── main.go                     # 程序入口
├── cmd/                        # Cobra 子命令
│   ├── gen.go                  # cpa gen：打包上下文
│   └── tree.go                 # cpa tree：打印文件树
├── internal/
│   ├── config/                 # TOML 配置解析
│   ├── gitignore/              # .gitignore 规则引擎
│   ├── tree/                   # 文件树生成逻辑
│   └── content/                # Markdown 内容构建
├── assets/
│   ├── assets.go               # Go embed 加载
│   └── skel/                   # 模板文件（嵌入到二进制中）
│       ├── config.toml         # 默认配置模板
│       ├── prompt.tmpl         # Go template 渲染模板
│       └── prompt_default.md   # 默认 AI 提示词
├── build/
│   └── .goreleaser.yml         # Goreleaser 配置
└── .github/workflows/          # CI/CD 流水线
    ├── ci.yml                  # 测试 + 构建
    ├── goreportcard.yml        # 代码质量检查
    └── release.yml             # 自动发布
```

---

## 技术栈（非常朴素）

| 技术                                                   | 用途            | 备注                |
| ------------------------------------------------------ | --------------- | ------------------- |
| Go 1.26+                                               | 主语言          | 我们爱 Go           |
| [Cobra](https://github.com/spf13/cobra)                 | CLI 框架        | 生态最成熟          |
| [go-toml](https://github.com/pelletier/go-toml)         | TOML 解析       | 比 V1 性能好 3 倍   |
| [embed](https://pkg.go.dev/embed)                       | 资源嵌入        | Go 官方包，零依赖   |
| [filepath.WalkDir](https://pkg.go.dev/filepath#WalkDir) | 文件遍历        | Go 1.16+ 的高效 API |
| [go-git/go-git](https://github.com/go-git/go-git)       | .gitignore 解析 | 纯 Go 实现          |

**总依赖：3 个直接依赖 + 5 个间接依赖。** 非常克制。

---

## 与其他工具的对比

| 对比项             | 手动开 N 个编辑器 | 其他 Prompt 工具 | **cpa**             |
| ------------------ | ----------------- | ---------------- | ------------------------- |
| 操作步骤           | 开 N 个 tab       | 配置复杂到想报警 | **一条命令 + 拖拽** |
| 支持 .gitignore    | ❌                | ❌               | ✅                        |
| 文件别名           | ❌                | ❌               | ✅                        |
| 自定义提示词       | ❌                | 可能支持         | ✅                        |
| 需要付费           | 电费              | 是               | **否**              |
| 需要网络请求       | 否                | 是               | **否**              |
| AI 兼容性          | N/A               | 部分             | **全宇宙**          |
| 装逼指数           | 😭                | 😐               | **😎**              |
| 被 AI 反客为主风险 | 低                | 低               | **极高** 🌀         |

---

## 贡献

欢迎提 Issue 和 PR！

```bash
# 本地开发
git clone https://github.com/yuan-shuo/copy-paste-ai.git
cd copy-paste-ai

# 快速验证
go run . tree
go run . gen -f main.go
go test ./... -v
```

**提交 PR 前请确保：**

- [ ] `go build ./...` 通过
- [ ] `go test ./...` 全部通过
- [ ] 圈复杂度 ≤ 15（`gocyclo -over 15 .`）

---

## License

自己看根目录

> 随便用，随便改，改坏了别找我，改好了记得给我点个 star 🌟

---

**Made with ❤️ and excessive usage of Ctrl+C / Ctrl+V**

**这个项目的存在证明了：当人类不想做一件事时，就会发明工具来逃避它。而现在我们又用 AI 来逃避写代码本身。**

**人类，真有趣。** 🤔
