package instruction

import (
	"context"
	"fmt"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/workflow"
	"log"
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
	workflow.SendBotMsgTemplate(
		ctx,
		l.Executor,
		e,
		l.createBotMsg,
	)
}

func (l InstructionInboundMsgListener) Listen(event any) {
	ctx := context.Background()
	workflow.ReceiveUserMsgTemplate(
		ctx,
		l.Executor,
		event,
		l.processUserMsg,
	)
}

func (l *InstructionStateListener) createBotMsg(
	ctx context.Context,
	event workflow.StateEvent,
	exec models.WorkflowExecution,
) (body string, action models.WorkflowAction, payload any) {
	switch event.StateAfter {
	case StateToSendOverview:
		body = fmt.Sprintf("%s\nReply with 1-7 to see details, 0 for overview again, negative number to quit.", loadReadMe()[0])
		action = ActionBotSentOverview
	case StateToSendSection:
		section, ok := event.ActionPayload.(int)
		if !ok {
			log.Printf("wrong action payload %v", event.ActionPayload)
			return
		}
		body = fmt.Sprintf("This is section %d. \n\n%s\nReply with 1-7 to see details, 0 for overview again, negative number to quit.", section, loadReadMe()[section])
		action = ActionBotSentSection
	case StateToSendGoobye:
		body = "Thank you for your time! Feel free to reach out!"
		action = ActionBotSentGoodbye
	default:
	}
	return body, action, payload
}

func (l *InstructionInboundMsgListener) processUserMsg(
	ctx context.Context,
	msg workflow.InboundMsgEvent,
	exec models.WorkflowExecution,
) (action models.WorkflowAction, actionPayload any) {
	switch exec.CurrentState {
	case StateWaitingForSelection:
		r, ok := l.NlpService.GetIntResult(ctx, msg.Body)
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
	}
	return
}
