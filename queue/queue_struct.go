/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 13:55:46
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-16 17:34:55
 * @Description: file content
 */
package queue

import "sync"

// 一个执行队列中的任务
type Task struct {
	ID       string
	Action   func() (map[string]interface{}, error)
	Callback func()
}

// 枚举事件类型
const (
	// 任务执行成功
	EventSuccess = 0
	// 任务执行失败
	EventFaile = 1
	// 任务被取消
	EventCancel = 2
	// 任务排队中
	EventPendding = 3
	// 任务执行中
	EventRunning = 4
)

// 事件消息
type EventMessage struct {
	ID        string                 `json:"id"`
	EventType int                    `json:"event_type"`
	EventData map[string]interface{} `json:"event_data"`
	Aq        *ActionQueue
}

// 一个执行队列
type ActionQueue struct {
	// 最大并发数
	MaxConcurrent int
	// 当前并发数
	CurrentConcurrent int
	// 任务队列
	TaskQueue chan Task
	// 事件队列
	EventQueue chan EventMessage
	// 任务ID列表
	TaskList []string
	// 互斥锁，保证并发安全
	mutex sync.Mutex
}
