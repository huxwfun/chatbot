package review

import (
	"huxwfun/chatbot/internal/models"
)

var ReviewWorkflow = models.Workflow{
	Name:               "review",
	Graph:              CreateReviewWorkflowGraph(),
	StateListener:      "review.ReviewStateListener",
	InboundMsgListener: "review.ReviewInboundMsgListener",
}
