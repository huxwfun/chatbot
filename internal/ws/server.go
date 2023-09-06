package ws

import (
	"context"
	"encoding/json"
	"flag"
	"huxwfun/chatbot/internal/event"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const BOT_CHAT_INBOUND_MSG = "BOT_CHAT_INBOUNT_MSG"
const BOT_CHAT_OUTBOUND_MSG = "BOT_CHAT_OUTBOUNT_MSG"

var addr = flag.String("addr", ":8080", "http service address")

type WsServer struct {
	dispatcher   *event.Dispatcher
	storage      *storage.Storage
	handleReview func(http.ResponseWriter, *http.Request)
	handleLog    func(http.ResponseWriter, *http.Request)
}

func NewWsServer(
	dispatcher *event.Dispatcher,
	storage *storage.Storage,
	handleReview func(http.ResponseWriter, *http.Request),
	handleLog func(http.ResponseWriter, *http.Request),
) *WsServer {
	return &WsServer{
		dispatcher,
		storage,
		handleReview,
		handleLog,
	}
}

func (ws *WsServer) sendMsg(ctx context.Context, msg models.Message) {
	chat, ok := ws.storage.Chat.Get(ctx, msg.ChatId)
	if !ok {
		log.Printf("error chat:%s missing", msg.ChatId)
		return
	}
	if online, ok := clients[chat.CustomerId]; ok && len(online) > 0 {
		for _, c := range online {
			bytes, err := json.Marshal(msg)
			if err != nil {
				log.Print(err)
				return
			}
			c.send <- bytes
		}
	}
}

func (ws *WsServer) listenForOutbound() {
	ws.dispatcher.Register(func(event interface{}) {
		ctx := context.Background()
		msg, ok := event.(models.Message)
		if !ok {
			log.Printf("error wrong message type")
			return
		}
		ws.sendMsg(ctx, msg)
	}, BOT_CHAT_OUTBOUND_MSG)
}

func (ws *WsServer) forwardInbound() {
	ws.dispatcher.Register(func(event interface{}) {
		ctx := context.Background()
		msg, ok := event.(models.Message)
		if !ok {
			return
		}
		ws.sendMsg(ctx, msg)
	}, BOT_CHAT_INBOUND_MSG)
}

func (ws *WsServer) Start() {
	flag.Parse()
	fs := http.FileServer(http.Dir("./frontend/out"))
	http.Handle("/", fs)
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		users := ws.storage.User.GetAll(ctx)
		var me models.User
		for {
			me = users[rand.Intn(len(users))]
			if !me.IsBot {
				break
			}
		}
		chats := ws.storage.Chat.FindByUser(ctx, me.Id)
		messages := ws.storage.Chat.FindMessagesByUser(ctx, me.Id)
		data := struct {
			Users    []models.User    `json:"users"`
			Me       models.User      `json:"me"`
			Chats    []models.BotChat `json:"chats"`
			Messages []models.Message `json:"messages"`
		}{
			Users:    users,
			Me:       me,
			Chats:    chats,
			Messages: messages,
		}
		body, err := json.Marshal(data)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err != nil {
			w.Header().Add(r.Response.Status, "500")
		} else {
			w.Write(body)
		}
	})
	http.HandleFunc("/review", ws.handleReview)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, ws.dispatcher)
	})
	http.HandleFunc("/logs", ws.handleLog)
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	ws.forwardInbound()
	ws.listenForOutbound()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}
