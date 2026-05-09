# 测试报告: `config` 模块

## 概述
`config` 模块提供了一种灵活的方式，可以从文件（JSON/YAML）加载配置，并使用环境变量进行覆盖。

## 测试场景
- **复杂结构映射**: 验证嵌套结构、映射 (Map)、指针以及 `Duration` 类型。
- **环境变量覆盖**: 确保配置能被环境变量正确覆盖，包括深层嵌套的映射。
- **回归测试**:
  - `TestForMap_Regression`: 验证基础的基于 Map 的配置加载。
  - `TestForStruct_Regression`: 验证基于结构体的映射及环境变量覆盖。
  - `TestForYml_Regression`: 检查 YAML 配置解析，包括接口切片和映射。
  - `TestForYaml_Support`: 验证针对无后缀 `.yaml` 文件的默认解析支持。

## 基准测试结果

```text
goos: darwin
goarch: amd64
pkg: apigo.cc/go/config
BenchmarkLoad-16                            6255            164111 ns/op
BenchmarkApplyEnvOverrides-16             294292              3592 ns/op
```

*注：基准测试反映了在 Intel i9-9980HK 上的性能表现。*
