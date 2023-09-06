package review

import (
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/workflow"
)

var ReviewWorkflow = models.Workflow{
	Name:  "review",
	Graph: CreateReviewWorkflowGraph(),
}

func CreateReviewWorkflow(
	id string,
	executor *workflow.WorkflowExecutor,
	nlpService *nlp.NlpService,
) models.Workflow {
	ReviewWorkflow.Id = id
	RegisterReviewStateListener(id, executor, nlpService)
	RegisterInboundMsgListener(id, executor, nlpService)
	return ReviewWorkflow
}
