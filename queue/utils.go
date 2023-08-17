/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 13:57:27
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-16 13:57:31
 * @Description: file content
 */
package queue

// 检查任务是否在在需要执行的列表中
func (aq *ActionQueue) checkTask(id string) bool {
	for _, taskID := range aq.TaskList {
		if taskID == id {
			return true
		}
	}
	return false
}
