# 扩展网络层协议

由于无法预估实际使用场景，所以可能本仓库提供的基础 `Point` 实现并不满足实际使用场景。
因此可以在 `pkg/plugins/point` 目录下完全重新按照 `pkg/plugins/point/point.go` 实现一套新方法。
然后重新注册到 `point` 中。

## 1. Point 说明

本章节将详细说明 `Point` 的相关类型方法定义说明。

### 1.1 Point 接口
```go
// source code in pkg/plugins/point/point.go 

package point

// Point define the protocol of TransLayout of RPC.
type Point interface {
    Name() string
    Client(id string, config plugins.AnyConfig, logger *logrus.Logger) (Client, error)
    Server(id string, config plugins.AnyConfig, logger *logrus.Logger) (Server, error)
}


// RegisterPoint will append or cover a point instance to inner point set. GetPoint will get point from this set.
func RegisterPoint(instance Point) {
	points[instance.Name()] = instance
}

// GetPoint will return specify point instance from inner set by name. And not found specify point, return false.
func GetPoint(name string) (Point, bool) {
	res, ok := points[name]
	return res, ok
}

```

`Point` 的目的是创建 `Client` 和 `Server` 实例所抽象的一个结构。
`Point` 本身**不应该**包含任何配置。如果具有单独的配置，请确保在传入 `RegisterPoint(instance)` 中的 `instance` 已经完成初始化动作。
