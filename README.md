## Overview
1. How to use
2. Workflow definition, execution and events log
3. Websocket and http interface
4. Processor chains for user input processing
5. In memory event bus for dispatching: workflow execution event, websocket messaging event
6. Simple UI built by React, Next.js 
7. Grossories

## How to use
docker run -p 8080:8080 $(docker build -q .)

http://localhost:8080

Press "Trigger sample review flow!" button on the right side to start bot reviewing.
## Workflow
[source code link](https://github.com/huxwfun/chatbot/tree/main/internal/workflow)

Separation of concerns:
- Define workflow using directed cyclic graph
- Use WorkflowExecution to represent the execution state of workflow, per customer
- Using an Executor to run workflow executions

Take review workflow for example:

1. Create the graph, each node is a state, each edge is a possible action when the workflow is in that state. Register listeners to the Executor
2. There shall be an entry node with state "WorkflowStateStart" and an edge to the start point of real business.
3. When some action triggers the start of "review" workflow, the Executor creates an WorkflowExecution object in state "WorkflowStateStart", with the customer and bot binded, and save it to DB
4. Executor does action "begin" to the execution, put it's state to "to prompt review"
5. ReviewStatelistener catches the coresponding event, sends the first message to user, then initiates another action "botPromptedReview" by Executor to the Execution. 
6. Executor puts Execution to state "waiting for confirmation", action event sent to ReviewStatelistener again, but this event won't trigger any further operations.
7. Nothing shall happen until the user answers the question. The ReviewInboundMsgListener catches user's message, and triggers next action based on analysis result of that message.

## Websocket
Copied from gorilla/websocket examples.

Basiccaly, once any client connects, the server creates a "client" instance, which uses 2 goroutines, 1 for reading, 1 for writing.

The reading goroutine sends msg to event bus directly

The writing goroutine watches a channel for msg to be sent to client on the other side of the websocket

Server registers another listener to event bus. Messages catched by this listener, classified by clients, and send to clients' writing goroutines through their channels

## User input processing

Rather simple processing chain.

## In memory event bus

Copied somewhere as a mock of a real event queue with delivery guarantees.

## UI
Typescript

React.js

Material UI

Next.js

## Grossories
## [Models](https://github.com/huxwfun/chatbot/tree/main/internal/models)
- Chat (represents user-bot conversation)
- Message (messages sent in chats)
- User (both users and bots)
- Workflow (workflow definition)
- WorkflowExecution (execution of workflow per customer)
- WorkflowLog (domain events during execution)
