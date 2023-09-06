package review

import (
	"context"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/utils"
	"huxwfun/chatbot/internal/workflow"
	"log"
	"time"
)

type ReviewStateListener struct {
	Executor   *workflow.WorkflowExecutor
	NlpService *nlp.NlpService
}

type ReviewInboundMsgListener struct {
	workflow.StateListener
	Executor   *workflow.WorkflowExecutor
	NlpService *nlp.NlpService
}

func (l ReviewStateListener) Listen(e any) {
	ctx := context.Background()
	event, ok := e.(workflow.StateEvent)
	if !ok {
		log.Printf("wrong event type")
		return
	}
	var (
		action = event.Action
		state  = event.StateAfter
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
	case StateToPromptReview:
		msg.Body = "Hello again! We noticed you've recently received your iPhone 13. We'd love to hear about your experience. Can you spare a few minutes to share your thoughts?"
		action = ActionBotPromptedReview
	case StateToAskForRating:
		msg.Body = "Fantastic! On a scale of 1-5, how would you rate the iPhone 13?"
		action = ActionBotAskedForRating
	case StateToSendCompletedGoobye:
		msg.Body = "Thank you for sharing your feedback! If you have any more thoughts or need assistance with anything else, feel free to reach out!"
		action = ActionBotSaidGoodbyte
	case StateToSendInterruptedGoobye:
		msg.Body = "Thank you for time! See you soon."
		action = ActionBotSaidGoodbyte
	default:
		return
	}
	l.Executor.SendBotMsg(ctx, msg)
	l.Executor.Action(ctx, exec, msg.Id, action, nil)
}

func (l ReviewInboundMsgListener) Listen(event any) {
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
	case StateWaitingForConfirmation:
		r, ok := l.NlpService.GetBoolResult(ctx, body)
		if !ok {
			action = ActionUnknown
		} else if r {
			action = ActionUserConfirmed
			actionPayload = true
		} else {
			action = ActionUserRejected
			actionPayload = false
		}
	case StateWaitingForRating:
		r, ok := l.NlpService.GetIntResult(ctx, body)
		if !ok {
			action = ActionUnknown
		} else {
			action = ActionUserSetRating
			actionPayload = r
		}
	case StateToPromptReview:
	case StateToAskForRating:
	case StateToSendCompletedGoobye:
	case StateToSendInterruptedGoobye:
	case StateDone:
	}

	l.Executor.Action(ctx, exec, msg.Id, action, actionPayload)
}
