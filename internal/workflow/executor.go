package workflow

import (
	"context"
	"huxwfun/chatbot/internal/event"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/storage"
	"huxwfun/chatbot/internal/utils"
	"huxwfun/chatbot/internal/ws"
	"log"
	"time"
)

type WorkflowExecutor struct {
	Storage    *storage.Storage
	Dispatcher *event.Dispatcher
}

func NewWorkflowExecutor(
	dispatcher *event.Dispatcher,
	storage *storage.Storage,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		Storage:    storage,
		Dispatcher: dispatcher,
	}
}

func (e *WorkflowExecutor) RegisterStateListener(workflowId string, l StateListener) {
	e.Dispatcher.Register(l.Listen, generateStateEventName(workflowId))
}

func (e *WorkflowExecutor) RegisterInboundMsgListener(workflowId string, l InboundMsgListener) {
	e.Dispatcher.Register(l.Listen, generateInboundMsgEventName(workflowId))
}

func (e *WorkflowExecutor) listenForInboundMsg() {
	e.Dispatcher.Register(func(event interface{}) {
		if msg, ok := event.(models.Message); ok {
			ctx := context.Background()
			e.Storage.Chat.SaveMessage(ctx, msg)

			executions := e.Storage.WorkflowExecution.FindActiveByCustomer(ctx, msg.From)
			for _, exec := range executions {
				e.Dispatcher.Dispatch(generateInboundMsgEventName(exec.WorkflowId), InboundMsgEvent{Message: msg, ExecutionId: exec.Id})
			}
		}
	}, ws.BOT_CHAT_INBOUND_MSG)
}

func (e *WorkflowExecutor) Action(
	ctx context.Context,
	exec models.WorkflowExecution,
	msgId string,
	action string,
	actionPayload interface{}) {
	workflow, ok := e.Storage.Workflow.Get(ctx, exec.WorkflowId)
	if !ok {
		log.Printf("error workflow:%s missing", exec.WorkflowId)
		return
	}
	before := exec.CurrentState
	edge := workflow.Graph.GetEdge(before, action)
	if edge == nil {
		log.Printf("error action %s cannot happen on state: %s", action, before)
		return
	}
	exec.CurrentState = edge.To
	log := models.WorkflowLog{
		Id:            utils.GenerateId(),
		ExecutionId:   exec.Id,
		StateBefore:   before,
		StateAfter:    exec.CurrentState,
		Action:        action,
		ActionPayload: actionPayload,
		MessageId:     msgId,
		TimeCreated:   time.Now(),
	}
	e.Storage.WorkflowExecution.Save(ctx, exec.Id, exec)
	e.Storage.WorkflowLog.Save(ctx, log.Id, log)
	e.Dispatcher.Dispatch(generateStateEventName(workflow.Id), log)
}

func (e *WorkflowExecutor) SendBotMsg(
	ctx context.Context, msg models.Message) {
	e.Dispatcher.Dispatch(ws.BOT_CHAT_OUTBOUND_MSG, msg)
	e.Storage.Chat.SaveMessage(ctx, msg)
}

func (e *WorkflowExecutor) CreateExecution(ctx context.Context, workflowId, userId, botId, chatId string) string {
	executions := e.Storage.WorkflowExecution.FindByCustomer(ctx, userId)
	for _, exec := range executions {
		if exec.CurrentState != models.WorkflowStateDone {
			e.InterruptExecution(ctx, exec)
		}
	}
	if workflow, ok := e.Storage.Workflow.Get(ctx, workflowId); ok {
		exec := models.WorkflowExecution{
			Id:           utils.GenerateId(),
			WorkflowId:   workflow.Id,
			CurrentState: models.WorkflowStateStart,
			BotId:        botId,
			CustomerId:   userId,
			ChatId:       chatId,
		}
		e.Storage.WorkflowExecution.Save(ctx, exec.Id, exec)
		e.Action(ctx, exec, "", models.WorkflowActionBegin, nil)
		return exec.Id
	}
	return ""
}

func (e *WorkflowExecutor) InterruptExecution(ctx context.Context, exec models.WorkflowExecution) {
	before := exec.CurrentState
	exec.CurrentState = models.WorkflowStateDone
	log := models.WorkflowLog{
		Id:            utils.GenerateId(),
		ExecutionId:   exec.Id,
		StateBefore:   before,
		StateAfter:    exec.CurrentState,
		Action:        models.WorkflowActionInterrupt,
		ActionPayload: nil,
		MessageId:     "",
		TimeCreated:   time.Now(),
	}
	e.Storage.WorkflowExecution.Save(ctx, exec.Id, exec)
	e.Storage.WorkflowLog.Save(ctx, log.Id, log)
}

func (e *WorkflowExecutor) Start() {
	e.listenForInboundMsg()
}
