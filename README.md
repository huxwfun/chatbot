## Workflow

Separation of concerns:
- Define workflow using directed cyclic graph
- Use WorkflowExecution to represent the execution state of workflow, per customer
- Using an Executor to run workflow executions


## Models

- Chat (represents user-bot conversation)
- Message (messages sent in chats)
- User (both users and bots)
- Workflow (workflow definition)
- WorkflowExecution (execution of workflow per customer)
- WorkflowLog (domain events during execution)
