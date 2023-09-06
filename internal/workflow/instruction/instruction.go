package instruction

import (
	"huxwfun/chatbot/internal/models"
)

var InstructionWorkflow = models.Workflow{
	Name:               "instruction",
	Graph:              CreateInstructionWorkflowGraph(),
	StateListener:      "instruction.InstructionStateListener",
	InboundMsgListener: "instruction.InstructionInboundMsgListener",
}
