/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 13:48:01
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-16 17:31:19
 * @Description: file content
 */
package queue

// 向队列放入事件
func (aq *ActionQueue) addEvent(id string, eventType int, eventData map[string]interface{}) {
	aq.EventQueue <- EventMessage{
		ID:        id,
		EventType: eventType,
		EventData: eventData,
		Aq:        aq,
	}
}

// 向队列放入错误事件
func (aq *ActionQueue) addErrorEvent(id string, reason string) {
	aq.addEvent(id, EventFaile, map[string]interface{}{
		"reason": reason,
	})
}

// 向队列放入排队事件
func (aq *ActionQueue) addPenddingEvent(id string) {
	aq.addEvent(id, EventPendding, nil)
}

// 向队列放入执行中事件
func (aq *ActionQueue) addRunningEvent(id string) {
	aq.addEvent(id, EventRunning, nil)
}

// 向队列放入成功事件
func (aq *ActionQueue) addSuccessEvent(id string, result map[string]interface{}) {
	aq.addEvent(id, EventSuccess, result)
}
