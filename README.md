## [Workflow](https://github.com/huxwfun/chatbot/tree/main/internal/workflow)

Separation of concerns:
- Define workflow using directed cyclic graph
- Use WorkflowExecution to represent the execution state of workflow, per customer
- Using an Executor to run workflow executions

## Workflow Execution
[CreateReiewWorkflow](https://github.com/huxwfun/chatbot/blob/main/internal/workflow/review/review.go)
- 

## [Models](https://github.com/huxwfun/chatbot/tree/main/internal/models)
- Chat (represents user-bot conversation)
- Message (messages sent in chats)
- User (both users and bots)
- Workflow (workflow definition)
- WorkflowExecution (execution of workflow per customer)
- WorkflowLog (domain events during execution)
