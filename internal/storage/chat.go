package storage

import (
	"context"
	"huxwfun/chatbot/internal/models"
	"huxwfun/chatbot/internal/utils"
)

type ChatStorage struct {
	InMemStorage[models.BotChat]
	messageStorage InMemStorage[models.Message]
}

func NewChatStorage() *ChatStorage {
	inMem := NewInMemStorage[models.BotChat]()
	messageStorage := NewInMemStorage[models.Message]()
	return &ChatStorage{
		inMem,
		messageStorage,
	}
}

func (c *ChatStorage) FindByUser(ctx context.Context, userId string) []models.BotChat {
	chats := []models.BotChat{}
	for _, v := range c.storage {
		if v.CustomerId == userId {
			chats = append(chats, v)
		}
	}
	return chats
}

func (c *ChatStorage) FindByUserAndBot(ctx context.Context, userId string, botId string) models.BotChat {
	for _, v := range c.storage {
		if v.CustomerId == userId && v.BotId == botId {
			return v
		}
	}
	chat := models.BotChat{
		Id:         utils.GenerateId(),
		CustomerId: userId,
		BotId:      botId,
	}
	c.Save(ctx, chat.Id, chat)
	return chat
}

func (c *ChatStorage) SaveMessage(ctx context.Context, msg models.Message) {
	c.messageStorage.Save(ctx, msg.Id, msg)
}

func (c *ChatStorage) GetMessage(ctx context.Context, msgId string) (models.Message, bool) {
	return c.messageStorage.Get(ctx, msgId)
}
func (c *ChatStorage) FindMessagesByUser(ctx context.Context, userId string) []models.Message {
	chats := c.FindByUser(ctx, userId)
	chatIds := map[string]bool{}
	for _, chat := range chats {
		chatIds[chat.Id] = true
	}
	messages := []models.Message{}
	for _, v := range c.messageStorage.storage {
		if chatIds[v.ChatId] {
			messages = append(messages, v)
		}
	}
	return messages
}
