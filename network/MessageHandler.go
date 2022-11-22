package network

import (
	"fmt"
	"server/iface"
	"server/utils"
	"strconv"
)

type MessageHandler struct {
	// 存放每个MsgId对应的处理方法
	ApiMap map[uint32]iface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	// Worker池中worker的数量
	WorkerPoolSize uint32
}

func NewMessageHandler() iface.IMessageHandler {
	return &MessageHandler{
		ApiMap:         make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObj.WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObj.WorkerPoolSize),
	}
}

// DoHandle 执行request所绑定的Router的业务
func (m *MessageHandler) DoHandle(request iface.IRequest) {
	msgId := request.GetMsgId()
	if router, exist := m.ApiMap[msgId]; !exist {
		panic("[MessageHandler Handle Router ERROR] Message id = +" +
			strconv.Itoa(int(msgId)) + ", missing router")
	} else {
		// 调用router的模版方法
		iface.Handle(router, request)
	}
}

func (m *MessageHandler) AddRouter(msgId uint32, router iface.IRouter) {
	if _, exist := m.ApiMap[msgId]; !exist {
		m.ApiMap[msgId] = router
		fmt.Printf("[Server Register Router] Message %d, add router success\n", msgId)
	} else {
		panic("[MessageHandler Register Router WARNING] Repeat api, msgID = " + strconv.Itoa(int(msgId)))
	}
}

// StartWorkerPool 启动一个Worker工作池 (开启工作池的动作只能发生一次)
// 根据WorkerPoolSize 分别开启Worker，每个Worker都用一个go来承载
func (m *MessageHandler) StartWorkerPool() {
	fmt.Printf("[Message Handler Worker Pool] Worker pool starting...\n")
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 为task queue开辟空间
		m.TaskQueue[i] = make(chan iface.IRequest, utils.MaxWorkerPoolSize)
		// start a worker and bond with a task queue
		go m.startWorker(m.TaskQueue[i], i)
	}
	fmt.Printf("[Message Handler Worker Pool] Worker pool started...\n")
}

// startWorker 启动一个Worker工作流程
func (m *MessageHandler) startWorker(taskQueue chan iface.IRequest, workerId int) {
	fmt.Printf("[Message Handler Worker] Worker %d is started\n", workerId)
	// 阻塞等待
	for {
		select {
		case request := <-taskQueue:
			// 执行当前request对应的业务
			m.DoHandle(request)
		}
	}
}

func (m *MessageHandler) SubmitTask(request iface.IRequest) {
	// 调度算法 -> 轮询分配 find a worker to handle this request
	workerId := request.GetConnection().GetConnId() % m.WorkerPoolSize
	fmt.Printf("[Message Handler Submit Task] Submit task %d to worker %d, connection id = %d\n",
		request.GetMsgId(), workerId, request.GetConnection().GetConnId())
	// send request to task queue of this worker
	m.TaskQueue[workerId] <- request
}
