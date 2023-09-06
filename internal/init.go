package chatbot

import (
	"context"
	"fmt"
	"huxwfun/chatbot/internal/event"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/storage"
	"huxwfun/chatbot/internal/workflow"
	"huxwfun/chatbot/internal/workflow/review"
	"huxwfun/chatbot/internal/ws"
	"log"
	"net/http"
)

func InitDispatcher(ctx context.Context) *event.Dispatcher {
	dis := event.NewDispatcher()
	return dis
}

func InitNlpSerivce(ctx context.Context) *nlp.NlpService {
	service := nlp.NewNlpService()
	service.Use(nlp.StaticBoolProcessor{})
	service.Use(nlp.NewSentimentBoolProcessor())
	service.Use(nlp.StaticIntProcessor{})
	return service
}

func InitStorage(ctx context.Context) *storage.Storage {
	storage := storage.NewStorage()
	log.Printf("storage initiated")
	return storage
}

func InitData(
	ctx context.Context,
	storage *storage.Storage,
) models.User {
	for _, customer := range CUSTOMERS {
		storage.User.Save(ctx, customer.Id, customer)
	}
	storage.User.Save(ctx, CHATBOG.Id, CHATBOG)
	for _, chat := range CHATS {
		storage.Chat.Save(ctx, chat.Id, chat)
	}
	log.Printf("data initiated, 3 users, 1 bot and 3 chats")
	return CHATBOG
}

func InitWorkflowExecutor(ctx context.Context,
	dispatcher *event.Dispatcher,
	storage *storage.Storage) *workflow.WorkflowExecutor {
	executor := workflow.NewWorkflowExecutor(dispatcher, storage)
	return executor
}

func RegisterWorkflow(
	ctx context.Context,
	executor *workflow.WorkflowExecutor,
	nlpService *nlp.NlpService,
	storage *storage.Storage) models.Workflow {
	id := REVIEW_WORKFLOW_ID
	workflow := review.CreateReviewWorkflow(id, executor, nlpService)

	storage.Workflow.Save(ctx, id, workflow)
	return workflow
}

func InitWebsocketServer(
	ctx context.Context,
	dispatcher *event.Dispatcher,
	storage *storage.Storage,
	executor *workflow.WorkflowExecutor,
	reviewWorkflow models.Workflow,
	reviewBot models.User,
) *ws.WsServer {
	log.Printf("init %s %s", reviewWorkflow.Id, reviewBot.Id)
	handleReview := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userId := r.URL.Query().Get("authentication")
		chat := storage.Chat.FindByUserAndBot(ctx, userId, reviewBot.Id)
		executionId := executor.CreateExecution(ctx, reviewWorkflow.Id, userId, reviewBot.Id, chat.Id)
		log.Printf("start review execution %s", executionId)
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	handleLog := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		output := ""
		userId := r.URL.Query().Get("authentication")
		executions := storage.WorkflowExecution.FindByCustomer(ctx, userId)
		for _, exec := range executions {
			bot, _ := storage.User.Get(ctx, exec.BotId)
			user, _ := storage.User.Get(ctx, exec.CustomerId)
			output = output + fmt.Sprintf("----- start execution %s(bot) to %s(you)-----\n", bot.Name, user.Name)
			logs := storage.WorkflowLog.FindByExecution(ctx, exec.Id)
			for _, log := range logs {
				output = output + fmt.Sprintf("%s\n", log.TimeCreated.Format("15:04:05"))
				if len(log.MessageId) > 0 {
					msg, ok := storage.Chat.GetMessage(ctx, log.MessageId)
					if ok {
						if msg.From == user.Id {
							output = output + fmt.Sprintf("%s send \"%s\" to %s\n", user.Name, msg.Body, bot.Name)
						} else {
							output = output + fmt.Sprintf("%s send \"%s\" to %s\n", bot.Name, msg.Body, user.Name)
						}
					} else {
						output = output + fmt.Sprintf("msg (%s) is missing\n", msg.Id)
					}
				}
				output = output + fmt.Sprintf("from state(%s) to state(%s) action(%s, payload: %v)\n\n", log.StateBefore, log.StateAfter, log.Action, log.ActionPayload)
			}
			output = output + fmt.Sprintf("----- end execution current state(%s) -----\n\n", exec.CurrentState)

		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(output))
	}
	return ws.NewWsServer(dispatcher, storage, handleReview, handleLog)
}
