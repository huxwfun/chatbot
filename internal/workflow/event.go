package workflow

import (
	"fmt"
	"huxwfun/chatbot/internal/models"
)

const StateEventName = "workflows:%s:state-changed"
const InboundMsgEventName = "workflows:%s:inbound-msg"

type StateEvent = models.WorkflowLog
type InboundMsgEvent struct {
	models.Message
	ExecutionId string
}

func generateStateEventName(workflowId string) string {
	return fmt.Sprintf(StateEventName, workflowId)
}
func generateInboundMsgEventName(workflowId string) string {
	return fmt.Sprintf(InboundMsgEventName, workflowId)
}

type StateListener interface {
	Listen(event any)
}
type InboundMsgListener interface {
	Listen(event any)
}
