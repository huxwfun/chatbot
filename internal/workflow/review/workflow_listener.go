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

func RegisterReviewStateListener(
	id string,
	executor *workflow.WorkflowExecutor,
	nlpService *nlp.NlpService,
) {
	executor.RegisterStateListener(id, func(e interface{}) {
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
		exec, ok := executor.Storage.WorkflowExecution.Get(ctx, event.ExecutionId)
		// executor.
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
		executor.SendBotMsg(ctx, msg)
		executor.Action(ctx, exec, msg.Id, action, nil)
	})
}
