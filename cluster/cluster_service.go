/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 15:17:45
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-18 13:04:46
 * @Description: file content
 */
package cluster

import (
	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/queue"

	"log"

	webui "github.com/SpenserCai/sd-webui-go"
)

type ClusterService struct {
	Nodes                 []*ClusterNode
	NodeEventQueue        chan queue.EventMessage
	GlobalEventQueue      chan queue.EventMessage
	NodesTaskMap          map[string]*ClusterNode
	Cfg                   *config.Config
	PreProcessEventStatus bool
}

func NewClusterService(cfg *config.Config) *ClusterService {
	nodes := make([]*ClusterNode, 0)
	nodeEventQuere := make(chan queue.EventMessage, 1000)
	// 从cfg的SDWebUI的Servers中获取所有的节点信息
	for _, server := range cfg.SDWebUi.Servers {
		node := NewClusterNode(
			server.Name,
			queue.NewActionQueue(server.MaxConcurrent, nodeEventQuere),
			webui.NewStableDiffInterface(server.Host),
		)
		nodes = append(nodes, node)
	}
	return &ClusterService{
		Nodes:                 nodes,
		NodeEventQueue:        nodeEventQuere,
		GlobalEventQueue:      make(chan queue.EventMessage, 1000),
		NodesTaskMap:          make(map[string]*ClusterNode),
		Cfg:                   cfg,
		PreProcessEventStatus: false,
	}
}

func (c *ClusterService) Start() {
	// 启动节点事件监听处理
	go c.PreProcessEvent()
	// 启动所有节点
	for _, node := range c.Nodes {
		go node.ActionQueue.Run()
	}
}

func (c *ClusterService) getNodeWithActionQueue(aq *queue.ActionQueue) *ClusterNode {
	for _, node := range c.Nodes {
		if node.ActionQueue == aq {
			return node
		}
	}
	return nil
}

func (c *ClusterService) PreProcessEvent() {
	c.PreProcessEventStatus = true
	defer func() {
		log.Println("PreProcessEvent Done")
		c.PreProcessEventStatus = false
	}()

	for {
		event := <-c.NodeEventQueue
		if event.EventType == queue.EventPendding {
			c.NodesTaskMap[event.ID] = c.getNodeWithActionQueue(event.Aq)
		} else {
			delete(c.NodesTaskMap, event.ID)
		}
		c.GlobalEventQueue <- event
	}
}

func (c *ClusterService) GetEvent() queue.EventMessage {
	return <-c.GlobalEventQueue
}

func (c *ClusterService) GetNode(name string) *ClusterNode {
	for _, node := range c.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (c *ClusterService) CancelTask(id string) {
	if node, ok := c.NodesTaskMap[id]; ok {
		node.ActionQueue.CancelTask(id)
	}
}

func (c *ClusterService) GetNodeAuto() *ClusterNode {
	// 返回一个负载最小的节点，如果都一样，返回第一个
	min := c.Nodes[0]
	for _, node := range c.Nodes {
		// 优先判断队列中的任务数量
		if len(node.ActionQueue.TaskList) < len(min.ActionQueue.TaskList) {
			min = node
		} else if len(node.ActionQueue.TaskList) == len(min.ActionQueue.TaskList) {
			if node.ActionQueue.CurrentConcurrent < min.ActionQueue.CurrentConcurrent {
				min = node
			}
		}
	}
	return min
}
