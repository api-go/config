# @go/config

> **Maintainer Statement:** 本项目完全由 AI 维护。任何改动均遵循代码质量与性能的最佳实践。


## 🎯 设计哲学

`@go/config` 旨在为应用提供一种**零配置感**的配置加载体验。其核心目标是：

*   **零摩擦入口**：`To[T]` 支持一键加载配置对象；`Load` 支持注入式加载。
*   **多源融合**：自动合并基础配置、环境特定配置（`env.json/yaml`）、命令行 `--config` 覆盖文件以及系统环境变量。
*   **语义映射**：通过下划线分隔的环境变量（如 `DB_HOST`）自动智能穿透映射到结构体嵌套字段（`{"db":{"host":...}}`）。
*   **极致敏捷**：支持 JSON/YAML 格式，两者共享 `cast` 的大小写与分隔符无关字段映射，并内置 `Duration` 等增强类型。

## 📦 安装

```bash
go get github.com/api-go/config
```

## 💡 快速开始

```go
import "github.com/api-go/config"

type MyConfig struct {
    Port int
    DB   struct {
        Host string
    }
}

// 1. 万能加载 (To 一键获取配置对象)
conf := config.To[MyConfig]("app") // 加载 app.json 或 app.yml

// 2. 传统注入 (Load 注入到现有指针)
var settings MyConfig
config.Load(&settings, "app")
```

## 🛠 API 指南

### 核心 API

1.  **通用加载 (Frictionless)**
    * `To[T any](name string) T` —— 万能加载入口。自动在当前目录、执行文件目录及用户主目录搜索 `name.json/yml`。失败时返回类型零值。
    * `Load(target any, name string) error` —— 注入加载。将配置内容注入到 `target` 指针。

2.  **加载机制与优先级**
    * **基础文件**: 搜索并加载 `name.json` 或 `name.yml`。
    * **环境补充**: 若存在 `env.json` 或 `env.yml`，则加载并覆盖基础文件内容。
    * **命令行覆盖**: `--config=env.1.yaml` 或 `--config env.1.yaml`。可重复指定并按出现顺序合并，后面的文件覆盖前面的文件。
    * **环境变量覆盖**: 自动扫描系统环境变量，按 `_` 分隔符映射到结构体字段（不区分大小写）。
        * 示例: `APP_PORT=8080` -> `{"app": {"port": 8080}}`

配置优先级为基础配置、`env.json/yaml`、命令行 `--config` 文件、系统环境变量。`--config` 参数会被后台启动器原样传递，因此管理进程和真实服务进程使用相同覆盖配置。

3.  **增强类型**
    * `Duration` —— 支持字符串解析的 `time.Duration` 包装类（如 `"5s"`, `"1h"`），解决标准库 `time.Duration` 在 JSON/YAML 反序列化中的限制。

## 🧪 验证状态
测试全部通过，支持复杂嵌套结构与环境变量重载。

详见：[TEST.md](./TEST.md)
