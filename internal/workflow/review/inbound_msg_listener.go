package review

import (
	"context"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/workflow"
	"log"
)

func RegisterInboundMsgListener(
	id string,
	executor *workflow.WorkflowExecutor,
	nlpService *nlp.NlpService,
) {
	executor.RegisterInboundMsgListener(id, func(event interface{}) {
		ctx := context.Background()
		msg, ok := event.(models.Message)
		if !ok {
			log.Printf("wrong msg type")
			return
		}
		body := msg.Body
		customerId := msg.From
		executions := executor.Storage.WorkflowExecution.FindByWorkflowAndCustomer(ctx, id, customerId)
		for _, exec := range executions {
			var action models.WorkflowAction
			var actionPayload interface{} = nil
			switch exec.CurrentState {
			case StateWaitingForConfirmation:
				r, ok := nlpService.GetBoolResult(ctx, body)
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
				r, ok := nlpService.GetIntResult(ctx, body)
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
			executor.Action(ctx, exec, msg.Id, action, actionPayload)
		}
	})
}
