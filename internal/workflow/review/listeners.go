package review

import (
	"context"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/workflow"
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
	workflow.SendBotMsgTemplate(
		ctx,
		l.Executor,
		e,
		l.createBotMsg,
	)
}

func (l ReviewInboundMsgListener) Listen(event any) {
	ctx := context.Background()
	workflow.ReceiveUserMsgTemplate(
		ctx,
		l.Executor,
		event,
		l.processUserMsg,
	)
}

func (l *ReviewStateListener) createBotMsg(
	ctx context.Context,
	event workflow.StateEvent,
	exec models.WorkflowExecution,
) (body string, action models.WorkflowAction, payload any) {
	switch event.StateAfter {
	case StateToPromptReview:
		body = "Hello again! We noticed you've recently received your iPhone 13. We'd love to hear about your experience. Can you spare a few minutes to share your thoughts?"
		action = ActionBotPromptedReview
	case StateToAskForRating:
		body = "Fantastic! On a scale of 1-5, how would you rate the iPhone 13?"
		action = ActionBotAskedForRating
	case StateToSendCompletedGoobye:
		body = "Thank you for sharing your feedback! If you have any more thoughts or need assistance with anything else, feel free to reach out!"
		action = ActionBotSaidGoodbyte
	case StateToSendInterruptedGoobye:
		body = "Thank you for time! See you soon."
		action = ActionBotSaidGoodbyte
	default:
		return
	}
	return body, action, payload
}

func (l *ReviewInboundMsgListener) processUserMsg(
	ctx context.Context,
	msg workflow.InboundMsgEvent,
	exec models.WorkflowExecution,
) (action models.WorkflowAction, actionPayload any) {
	switch exec.CurrentState {
	case StateWaitingForConfirmation:
		r, ok := l.NlpService.GetBoolResult(ctx, msg.Body)
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
		r, ok := l.NlpService.GetIntResult(ctx, msg.Body)
		if !ok {
			action = ActionUnknown
		} else {
			action = ActionUserSetRating
			actionPayload = r
		}
	default:
	}
	return
}
