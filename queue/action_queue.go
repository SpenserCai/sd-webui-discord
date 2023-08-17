/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:15
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-16 17:05:59
 * @Description: file content
 */
package queue

import (
	"fmt"
	"log"
	"time"
)

// 创建一个执行队列
func NewActionQueue(maxConcurrent int, eventQueue chan EventMessage) *ActionQueue {
	// 如果eventQueue为空，则创建一个，否则使用传入的eventQueue
	if eventQueue == nil {
		eventQueue = make(chan EventMessage, 100)
	}
	return &ActionQueue{
		MaxConcurrent:     maxConcurrent,
		CurrentConcurrent: 0,
		TaskQueue:         make(chan Task, 100),
		EventQueue:        eventQueue,
		TaskList:          make([]string, 0),
	}
}

// 向执行队列中添加一个任务
func (aq *ActionQueue) AddTask(id string, task func() (map[string]interface{}, error), callback func()) {
	aq.mutex.Lock()
	defer aq.mutex.Unlock()

	if aq.checkTask(id) {
		aq.addErrorEvent(id, "task id is exist")
		return
	}

	log.Println("AddTask:", id)
	aq.TaskQueue <- Task{
		ID:       id,
		Action:   task,
		Callback: callback,
	}
	aq.TaskList = append(aq.TaskList, id)
	log.Println("AddTask success:", id)
	aq.addPenddingEvent(id)
}

// 根据任务ID取消任务
func (aq *ActionQueue) CancelTask(id string) {
	aq.mutex.Lock()
	defer aq.mutex.Unlock()

	log.Println("CancelTask:", id)
	// 判断任务是否在队列中
	for i, taskID := range aq.TaskList {
		if taskID == id {
			// 从任务队列中移除
			aq.TaskList = append(aq.TaskList[:i], aq.TaskList[i+1:]...)
			aq.addEvent(id, EventCancel, nil)
			return
		}
	}
	aq.addErrorEvent(id, "task id is not exist or task is running")

}

// 获取错误队列中的错误消息
func (aq *ActionQueue) GetEvent() EventMessage {
	return <-aq.EventQueue
}

// 执行队列中的任务
func (aq *ActionQueue) Run() {
	for {
		if aq.CurrentConcurrent < aq.MaxConcurrent && len(aq.TaskQueue) > 0 {
			aq.mutex.Lock()
			aq.CurrentConcurrent++
			task := <-aq.TaskQueue
			if !aq.checkTask(task.ID) {
				aq.mutex.Unlock()
				continue
			}
			for i, taskID := range aq.TaskList {
				if taskID == task.ID {
					// 从任务队列中移除
					aq.TaskList = append(aq.TaskList[:i], aq.TaskList[i+1:]...)
					break
				}
			}
			aq.addRunningEvent(task.ID)
			aq.mutex.Unlock()

			// 当前并发数，当前任务队列长度
			logOut := fmt.Sprintf("CurrentConcurrent: %d, TaskQueueLen: %d", aq.CurrentConcurrent, len(aq.TaskQueue))
			log.Println(logOut)

			go func() {
				result, err := task.Action()
				if err != nil {
					aq.addErrorEvent(task.ID, err.Error())
				} else {
					if task.Callback != nil && result != nil {
						task.Callback()
					}
					aq.addSuccessEvent(task.ID, result)
				}

				aq.mutex.Lock()
				aq.CurrentConcurrent--
				aq.mutex.Unlock()
			}()
		}
		// 防止CPU占用过高
		time.Sleep(time.Millisecond * 10)
	}
}
