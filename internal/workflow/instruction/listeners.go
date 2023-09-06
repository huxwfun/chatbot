package instruction

import (
	"context"
	"fmt"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/utils"
	"huxwfun/chatbot/internal/workflow"
	"log"
	"time"
)

type InstructionStateListener struct {
	Executor   *workflow.WorkflowExecutor
	NlpService *nlp.NlpService
}

type InstructionInboundMsgListener struct {
	workflow.StateListener
	Executor   *workflow.WorkflowExecutor
	NlpService *nlp.NlpService
}

func (l InstructionStateListener) Listen(e any) {
	ctx := context.Background()
	event, ok := e.(workflow.StateEvent)
	if !ok {
		log.Printf("wrong event type")
		return
	}
	var (
		action        = event.Action
		actionPayload = event.ActionPayload
		state         = event.StateAfter
	)
	exec, ok := l.Executor.Storage.WorkflowExecution.Get(ctx, event.ExecutionId)
	if !ok {
		log.Printf("execution is missing")
		return
	}
	msg := models.Message{
		Id:          utils.GenerateId(),
		ChatId:      exec.ChatId,
		From:        exec.BotId,
		TimeCreated: time.Now(),
	}
	switch state {
	case StateToSendOverview:
		msg.Body = fmt.Sprintf("%s\nReply with 1-7 to see details, 0 for overview again, negative number to quit.", loadReadMe()[0])
		action = ActionBotSentOverview
	case StateToSendSection:
		section, ok := actionPayload.(int)
		if !ok {
			log.Printf("wrong action payload %v", actionPayload)
			return
		}
		msg.Body = loadReadMe()[section]
		action = ActionBotSentSection
	case StateToSendGoobye:
		msg.Body = "Thank you for your time! Feel free to reach out!"
		action = ActionBotSentGoodbye
	default:
		return
	}
	l.Executor.SendBotMsg(ctx, msg)
	l.Executor.Action(ctx, exec, msg.Id, action, nil)
}

func (l InstructionInboundMsgListener) Listen(event any) {
	ctx := context.Background()
	msg, ok := event.(workflow.InboundMsgEvent)
	if !ok {
		log.Printf("wrong msg type")
		return
	}
	body := msg.Body
	exec, ok := l.Executor.Storage.WorkflowExecution.Get(ctx, msg.ExecutionId)
	if !ok {
		log.Printf("execution %s is missing", msg.ExecutionId)
		return
	}
	var action models.WorkflowAction
	var actionPayload interface{} = nil
	switch exec.CurrentState {
	case StateWaitingForSelection:
		r, ok := l.NlpService.GetIntResult(ctx, body)
		if !ok {
			action = ActionUnknown
		} else if r > 7 {
			action = ActionUnknown
			actionPayload = r
		} else if r < 0 {
			action = ActionUserChooseDone
			actionPayload = r
		} else {
			action = ActionUserSelected
			actionPayload = r
		}
	default:
		return
	}

	l.Executor.Action(ctx, exec, msg.Id, action, actionPayload)
}
