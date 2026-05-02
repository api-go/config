# 配置模块 (config)

`config` 模块负责从文件（JSON/YAML）加载配置，并支持通过环境变量进行覆盖。

## API 指南
### `LoadConfig(name string, conf interface{}) error`
从当前目录、可执行文件目录或用户主目录中加载 `name.json` 或 `name.yml`。随后，它会尝试加载 `env.json` 或 `env.yml` 并使用环境变量进行最终覆盖。

## 环境变量覆盖规则
环境变量使用 `_` 作为分隔符映射到配置结构中的嵌套字段。
示例：`DB_HOST` 将映射到 `{"db": {"host": ...}}`。
