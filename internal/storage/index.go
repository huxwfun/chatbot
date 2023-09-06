package storage

type Storage struct {
	Chat              *ChatStorage
	User              *UserStorage
	Workflow          *WorkflowStorage
	WorkflowExecution *WorkflowExecutionStorage
	WorkflowLog       *WorkflowLogStorage
}

func NewStorage() *Storage {
	Chat := NewChatStorage()
	User := NewUserStorage()
	Workflow := NewWorkflowStorage()
	WorkflowExecution := NewWorkflowExecutionStorage()
	WorkflowLog := NewWorkflowLogStorage()
	return &Storage{
		Chat,
		User,
		Workflow,
		WorkflowExecution,
		WorkflowLog,
	}
}
