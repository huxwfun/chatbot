package workflow

import (
	"context"
	"fmt"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/utils"
	"log"
	"time"
)

const StateEventName = "workflows:%s:state-changed"
const InboundMsgEventName = "workflows:%s:inbound-msg"

type StateEvent = models.WorkflowLog
type InboundMsgEvent struct {
	models.Message
	ExecutionId string
}

func generateStateEventName(workflowId string) string {
	return fmt.Sprintf(StateEventName, workflowId)
}
func generateInboundMsgEventName(workflowId string) string {
	return fmt.Sprintf(InboundMsgEventName, workflowId)
}

type StateListener interface {
	Listen(event any)
}
type InboundMsgListener interface {
	Listen(event any)
}

// template message pattern
func SendBotMsgTemplate(
	ctx context.Context,
	executor *WorkflowExecutor,
	e any,
	process func(ctx context.Context, event StateEvent, exec models.WorkflowExecution) (string, models.WorkflowAction, any),
) {
	event, ok := e.(StateEvent)
	if !ok {
		log.Printf("wrong event type %v", e)
	}
	exec, ok := executor.Storage.WorkflowExecution.Get(ctx, event.ExecutionId)
	if !ok {
		log.Printf("execution %s is missing", event.ExecutionId)
	}
	body, action, actionPayload := process(ctx, event, exec)
	if len(body) > 0 {
		msg := models.Message{
			Id:          utils.GenerateId(),
			ChatId:      exec.ChatId,
			From:        exec.BotId,
			TimeCreated: time.Now(),
		}
		msg.Body = body
		executor.SendBotMsg(ctx, msg)
		if len(action) > 0 {
			executor.Action(ctx, exec, msg.Id, action, actionPayload)
		}
	}
}

func ReceiveUserMsgTemplate(
	ctx context.Context,
	executor *WorkflowExecutor,
	e any,
	process func(ctx context.Context, msg InboundMsgEvent, exec models.WorkflowExecution) (action models.WorkflowAction, acitonPayload any),
) {
	msg, ok := e.(InboundMsgEvent)
	if !ok {
		log.Printf("wrong msg type")
		return
	}
	exec, ok := executor.Storage.WorkflowExecution.Get(ctx, msg.ExecutionId)
	if !ok {
		log.Printf("execution %s is missing", msg.ExecutionId)
		return
	}
	action, actionPayload := process(ctx, msg, exec)
	if len(action) > 0 {
		executor.Action(ctx, exec, msg.Id, action, actionPayload)
	}
}
