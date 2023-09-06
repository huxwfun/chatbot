package models

import "time"

type WorkflowState = string
type WorkflowAction = string

const WorkflowStateStart WorkflowState = ""
const WorkflowStateDone WorkflowState = "done"

const WorkflowActionBegin WorkflowAction = "begin"
const WorkflowActionInterrupt WorkflowAction = "interrupt"

type Workflow struct {
	Id    string
	Name  string
	Graph Graph
}

type WorkflowExecution struct {
	Id           string
	WorkflowId   string
	CurrentState WorkflowState
	BotId        string
	CustomerId   string
	ChatId       string
}

type WorkflowLog struct {
	Id            string
	ExecutionId   string
	MessageId     string
	StateBefore   WorkflowState
	StateAfter    WorkflowState
	Action        WorkflowAction
	ActionPayload interface{}
	TimeCreated   time.Time
}
