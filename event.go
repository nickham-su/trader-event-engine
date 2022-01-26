package event_engine

// 事件
type Event struct {
	EventType EventType   // 事件类型
	Data      interface{} // 事件数据
}
