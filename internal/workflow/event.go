package workflow

import (
	"fmt"
	"huxwfun/chatbot/internal/models"
)

const StateEventName = "workflows:%s:state-changed"
const InboundMsgEventName = "workflows:%s:inbound-msg"

type StateEvent = models.WorkflowLog
type InboundMsgEvent = models.Message

func generateStateEventName(workflowId string) string {
	return fmt.Sprintf(StateEventName, workflowId)
}
func generateInboundMsgEventName(workflowId string) string {
	return fmt.Sprintf(InboundMsgEventName, workflowId)
}
