package main

import (
	"context"
	chatbot "huxwfun/chatbot/internal"
)

func main() {
	ctx := context.Background()
	dis := chatbot.InitDispatcher(ctx)
	storage := chatbot.InitStorage(ctx)
	reviewBot := chatbot.InitData(ctx, storage)
	executor := chatbot.InitWorkflowExecutor(ctx, dis, storage)
	nlpService := chatbot.InitNlpSerivce(ctx)
	reviewWorkflow := chatbot.RegisterWorkflow(ctx, executor, nlpService, storage)
	server := chatbot.InitWebsocketServer(ctx, dis, storage, executor, reviewWorkflow, reviewBot)
	executor.Start()
	server.Start()
}
