# trace

`trace` 是 infrago 的模块包。

## 安装

```bash
go get github.com/infrago/trace@latest
```

## 最小接入

```go
package main

import (
    _ "github.com/infrago/trace"
    "github.com/infrago/infra"
)

func main() {
    infra.Run()
}
```

## 配置示例

```toml
[trace]
driver = "default"
```

## 公开 API（摘自源码）

- `func SpanValues(span Span, instance, flag string) map[string]Any`
- `func ResolveFields(raw Any, defaults map[string]string) map[string]string`
- `func (inst *Instance) Allow(span Span) bool`
- `func (inst *Instance) AllowWithFactor(span Span, factor float64) bool`
- `func (inst *Instance) Format(span Span) string`
- `func (d *defaultDriver) Connect(inst *Instance) (Connection, error)`
- `func (c *defaultConnection) Open() error  { return nil }`
- `func (c *defaultConnection) Close() error { return nil }`
- `func (c *defaultConnection) Write(spans ...Span) error`
- `func Begin(meta *infra.Meta, name string, attrs ...Map) *Handle`
- `func (h *Handle) End(results ...Any)`
- `func Emit(meta *infra.Meta, name string, status string, attrs ...Map)`
- `func Write(span Span)`
- `func RegisterDriver(name string, driver Driver)`
- `func RegisterConfig(name string, cfg Config)`
- `func RegisterConfigs(configs Configs)`
- `func Stats() Map`
- `func (m *Module) Register(name string, value Any)`
- `func (m *Module) RegisterDriver(name string, driver Driver)`
- `func (m *Module) RegisterConfig(name string, cfg Config)`
- `func (m *Module) RegisterConfigs(configs Configs)`
- `func (m *Module) Config(global Map)`
- `func (m *Module) Setup()`
- `func (m *Module) Open()`
- `func (m *Module) Start()`
- `func (m *Module) Stop()`
- `func (m *Module) Close()`
- `func (m *Module) Write(span Span)`
- `func (m *Module) Stats() Map`
- `func (m *Module) Begin(meta *infra.Meta, name string, attrs Map) infra.TraceSpan`
- `func (m *Module) Trace(meta *infra.Meta, name string, status string, attrs Map) error`

## 排错

- 模块未运行：确认空导入已存在
- driver 无效：确认驱动包已引入
- 配置不生效：检查配置段名是否为 `[trace]`
