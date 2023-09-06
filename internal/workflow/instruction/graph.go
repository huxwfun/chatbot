package instruction

import (
	"huxwfun/chatbot/internal/models"
)

const (
	StateToSendOverview      models.WorkflowState = "to send overview"
	StateToSendSection       models.WorkflowState = "to send section"
	StateWaitingForSelection models.WorkflowState = "waiting for selection"
	StateToSendGoobye        models.WorkflowState = "to send goodbye"
	StateDone                                     = models.WorkflowStateDone
)

const (
	ActionUnknown         models.WorkflowAction = "unknown"
	ActionBotSentOverview models.WorkflowAction = "botSentOverview"
	ActionBotSentSection  models.WorkflowAction = "botSentSection"
	ActionUserSelected    models.WorkflowAction = "userSelected"
	ActionUserChooseDone  models.WorkflowAction = "userChooseDone"
	ActionBotSentGoodbye  models.WorkflowAction = "botSentGoodbye"
)

func CreateInstructionWorkflowGraph() models.Graph {
	graph := models.CreateGraph()

	graph.AddNode(models.WorkflowStateStart)
	graph.AddNode(StateToSendOverview)
	graph.AddNode(StateWaitingForSelection)
	graph.AddNode(StateToSendSection)
	graph.AddNode(StateToSendGoobye)
	graph.AddNode(StateDone)

	graph.AddEdge(models.WorkflowStateStart, StateToSendOverview, models.WorkflowActionBegin)

	graph.AddEdge(StateToSendOverview, StateWaitingForSelection, ActionBotSentOverview)

	graph.AddEdge(StateWaitingForSelection, StateToSendSection, ActionUserSelected)
	graph.AddEdge(StateWaitingForSelection, StateToSendGoobye, ActionUserChooseDone)

	graph.AddEdge(StateToSendSection, StateWaitingForSelection, ActionBotSentSection)

	graph.AddEdge(StateToSendGoobye, StateDone, ActionBotSentGoodbye)
	return graph
}
