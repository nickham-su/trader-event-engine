package eventEngine

import (
	"github.com/nickham-su/go-queue"
)

func New(newTickChannel <-chan interface{}, doneChannel <-chan interface{}) *EventEngine {
	return &EventEngine{
		new(queue.Queue),
		make(map[EventType][]EventHandler),
		make(map[EventType][]EventHandler),
		newTickChannel,
		doneChannel,
		make(map[string]interface{}),
	}
}

// 事件引擎
type EventEngine struct {
	queue           *queue.Queue                 // 事件队列
	handlers        map[EventType][]EventHandler // 事件处理方法
	generalHandlers map[EventType][]EventHandler // 全局事件处理方法
	newTickChannel  <-chan interface{}           // 新tick通道
	doneChannel     <-chan interface{}           // 完成通道
	context         map[string]interface{}       // 上下文环境
}

// 注册事件处理方法
func (e *EventEngine) Register(eventType EventType, eh EventHandler) {
	e.handlers[eventType] = append(e.handlers[eventType], eh)
}

// 注册全局时间处理方法
func (e *EventEngine) RegisterGeneral(eventType EventType, eh EventHandler) {
	e.generalHandlers[eventType] = append(e.generalHandlers[eventType], eh)
}

// 添加事件
func (e *EventEngine) Put(event *Event) {
	e.queue.Push(event)
}

// 运行事件引擎
func (e *EventEngine) Run() {
	for {
		if ev, ok := e.queue.Pop(); ok {
			if event, ok2 := ev.(*Event); ok2 {
				e.process(event)
			}
		} else {
			select {
			case tick := <-e.newTickChannel:
				e.process(&Event{
					EventType: EtNewTick,
					Data:      tick,
				})
			case <-e.doneChannel:
				return
			}
		}
	}
}

// 处理事件
func (e *EventEngine) process(event *Event) {
	et := event.EventType

	for _, eh := range e.handlers[et] {
		eh.Handler(event, e)
	}

	for _, eh := range e.generalHandlers[et] {
		eh.Handler(event, e)
	}
}

func (e *EventEngine) SetContext(key string, value interface{}) {
	e.context[key] = value
}

func (e *EventEngine) GetContext(key string) interface{} {
	return e.context[key]
}
