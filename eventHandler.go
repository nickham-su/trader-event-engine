package eventEngine

// 事件处理器接口
type EventHandler interface {
	Handler(*Event, *EventEngine) // 事件处理方法
}
