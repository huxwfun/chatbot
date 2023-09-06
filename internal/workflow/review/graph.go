package review

import (
	"huxwfun/chatbot/internal/models"
)

const (
	StateToPromptReview          models.WorkflowState = "to prompt review"
	StateWaitingForConfirmation  models.WorkflowState = "waiting for confirmation"
	StateToAskForRating          models.WorkflowState = "to ask for rating"
	StateWaitingForRating        models.WorkflowState = "waiting for rating"
	StateToSendCompletedGoobye   models.WorkflowState = "to send completed goodbye"
	StateToSendInterruptedGoobye models.WorkflowState = "to send interrupted goodbye"
	StateDone                                         = models.WorkflowStateDone
)

const (
	ActionUnknown           models.WorkflowAction = "unknown"
	ActionBotPromptedReview models.WorkflowAction = "botPromptedReview"
	ActionUserConfirmed     models.WorkflowAction = "userConfirmed"
	ActionUserRejected      models.WorkflowAction = "userRejected"
	ActionBotAskedForRating models.WorkflowAction = "botAskedForRating"
	ActionUserSetRating     models.WorkflowAction = "userSetRating"
	ActionBotSaidGoodbyte   models.WorkflowAction = "botSaidGoodbyte"
)

func CreateReviewWorkflowGraph() models.Graph {
	graph := models.CreateGraph()
	graph.AddNode(models.WorkflowStateStart)
	graph.AddNode(StateToPromptReview)
	graph.AddNode(StateWaitingForConfirmation)
	graph.AddNode(StateToAskForRating)
	graph.AddNode(StateWaitingForRating)
	graph.AddNode(StateToSendCompletedGoobye)
	graph.AddNode(StateToSendInterruptedGoobye)
	graph.AddNode(StateDone)

	graph.AddEdge(models.WorkflowStateStart, StateToPromptReview, models.WorkflowActionBegin)

	graph.AddEdge(StateToPromptReview, StateWaitingForConfirmation, ActionBotPromptedReview)
	graph.AddEdge(StateToPromptReview, StateWaitingForConfirmation, ActionBotPromptedReview)

	graph.AddEdge(StateWaitingForConfirmation, StateToAskForRating, ActionUserConfirmed)
	graph.AddEdge(StateWaitingForConfirmation, StateToSendInterruptedGoobye, ActionUserRejected)
	graph.AddEdge(StateWaitingForConfirmation, StateToPromptReview, ActionUnknown)

	graph.AddEdge(StateToAskForRating, StateWaitingForRating, ActionBotAskedForRating)

	graph.AddEdge(StateWaitingForRating, StateToSendCompletedGoobye, ActionUserSetRating)
	graph.AddEdge(StateWaitingForRating, StateToSendInterruptedGoobye, ActionUserRejected)
	graph.AddEdge(StateWaitingForRating, StateToAskForRating, ActionUnknown)

	graph.AddEdge(StateToSendCompletedGoobye, StateDone, ActionBotSaidGoodbyte)
	graph.AddEdge(StateToSendInterruptedGoobye, StateDone, ActionBotSaidGoodbyte)

	return graph
}
