package chatbot

import (
	"context"
	"fmt"
	"huxwfun/chatbot/internal/event"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/nlp"
	"huxwfun/chatbot/internal/storage"
	"huxwfun/chatbot/internal/workflow"
	"huxwfun/chatbot/internal/workflow/instruction"
	"huxwfun/chatbot/internal/workflow/review"
	"huxwfun/chatbot/internal/ws"
	"log"
	"net/http"
	"time"
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
	{
		id := REVIEW_WORKFLOW_ID
		workflow := review.ReviewWorkflow
		workflow.Id = id
		storage.Workflow.Save(ctx, id, workflow)
	}
	{
		id := INSTRUCTION_WORKFLOW_ID
		workflow := instruction.InstructionWorkflow
		workflow.Id = id
		storage.Workflow.Save(ctx, id, workflow)
	}
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
	storage *storage.Storage) {
	workflows := storage.Workflow.GetAll(ctx)
	for _, wf := range workflows {
		// let's pretend there's DI system to create instance from string
		if wf.InboundMsgListener == "review.ReviewInboundMsgListener" {
			inboundMsgListener := review.ReviewInboundMsgListener{
				Executor:   executor,
				NlpService: nlpService,
			}
			executor.RegisterInboundMsgListener(wf.Id, inboundMsgListener)
		}
		if wf.StateListener == "review.ReviewStateListener" {
			stateListener := review.ReviewStateListener{
				Executor:   executor,
				NlpService: nlpService,
			}
			executor.RegisterStateListener(wf.Id, stateListener)
		}
		if wf.StateListener == "instruction.InstructionStateListener" {
			inboundMsgListener := instruction.InstructionStateListener{
				Executor:   executor,
				NlpService: nlpService,
			}
			executor.RegisterStateListener(wf.Id, inboundMsgListener)
		}
		if wf.InboundMsgListener == "instruction.InstructionInboundMsgListener" {
			stateListener := instruction.InstructionInboundMsgListener{
				Executor:   executor,
				NlpService: nlpService,
			}
			executor.RegisterInboundMsgListener(wf.Id, stateListener)
		}
	}
}

func InitWebsocketServer(
	ctx context.Context,
	dispatcher *event.Dispatcher,
	storage *storage.Storage,
	executor *workflow.WorkflowExecutor,
	reviewBot models.User,
) *ws.WsServer {
	handleReview := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userId := r.URL.Query().Get("authentication")
		chat := storage.Chat.FindByUserAndBot(ctx, userId, reviewBot.Id)
		executionId := executor.CreateExecution(ctx, REVIEW_WORKFLOW_ID, userId, reviewBot.Id, chat.Id)
		log.Printf("start review execution %s", executionId)
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	handleInstruction := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userId := r.URL.Query().Get("authentication")
		chat := storage.Chat.FindByUserAndBot(ctx, userId, reviewBot.Id)
		executionId := executor.CreateExecution(ctx, INSTRUCTION_WORKFLOW_ID, userId, reviewBot.Id, chat.Id)
		log.Printf("start instruction execution %s", executionId)
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	handleLog := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		output := ""
		userId := r.URL.Query().Get("authentication")
		executions := storage.WorkflowExecution.FindByCustomer(ctx, userId)
		for _, exec := range executions {
			workflow, _ := storage.Workflow.Get(ctx, exec.WorkflowId)
			bot, _ := storage.User.Get(ctx, exec.BotId)
			user, _ := storage.User.Get(ctx, exec.CustomerId)
			output = output + fmt.Sprintf("## %s(bot) starts %s to %s(you)-----\n", bot.Name, workflow.Name, user.Name)
			logs := storage.WorkflowLog.FindByExecution(ctx, exec.Id)
			for _, log := range logs {
				output = output + fmt.Sprintf("#### %s: \n", log.TimeCreated.Format("15:04:05"))
				output = output + fmt.Sprintf("action: ***%s***, payload: *%v*\n\n", log.Action, log.ActionPayload)
				output = output + fmt.Sprintf("state: ***%s*** --> ***%s***\n\n", log.StateBefore, log.StateAfter)
				if len(log.MessageId) > 0 {
					msg, ok := storage.Chat.GetMessage(ctx, log.MessageId)
					output = output + fmt.Sprintf("```\n")
					if ok {
						if msg.From == user.Id {
							output = output + fmt.Sprintf("%s to %s:\n", user.Name, bot.Name)
						} else {
							output = output + fmt.Sprintf("%s to %s:\n", bot.Name, user.Name)
						}
						output = output + fmt.Sprintf("%s\n", msg.Body)
					} else {
						output = output + fmt.Sprintf("msg (%s) is missing\n", msg.Id)
					}
					output = output + fmt.Sprintf("```\n")
				}
			}
			output = output + fmt.Sprintf("#### %s current state ***%s***\n", time.Now().Format("15:04:05"), exec.CurrentState)

		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(output))
	}
	return ws.NewWsServer(dispatcher, storage, handleReview, handleInstruction, handleLog)
}
