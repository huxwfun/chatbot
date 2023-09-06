package chatbot

import "huxwfun/chatbot/internal/models"

var CUSTOMERS = []models.User{
	{
		Id:     "e8f6de79-6499-4f2b-8a25-199a37faca1c",
		Name:   "Sarah",
		Avatar: "1.avif",
	},
	{
		Id:     "1bd9c528-1887-499f-a25f-b99f8d7e9b82",
		Name:   "May",
		Avatar: "2.avif",
	},
	{
		Id:     "2f52406b-5dbc-40d0-bbcb-6b739dcecc5f",
		Name:   "Alice",
		Avatar: "3.avif",
	},
}

var CHATBOG = models.User{
	Id:     "fa935e1d-3124-437d-bcc4-7349c4b95189",
	Name:   "R2D20",
	Avatar: "bot.jpeg",
	IsBot:  true,
}

var REVIEW_WORKFLOW_ID = "86fe916a-531b-4216-bc7e-69122c6d91a4"

var CHATS = []models.BotChat{
	{
		Id:         "6cd26eb5-dc0b-4c28-b1dc-33f5389e6e41",
		BotId:      CHATBOG.Id,
		CustomerId: CUSTOMERS[0].Id,
	},
	{
		Id:         "0991ab62-b500-495e-8430-89a67655997b",
		BotId:      CHATBOG.Id,
		CustomerId: CUSTOMERS[1].Id,
	},
	{
		Id:         "ac0565ac-3a8f-471f-8ce7-c54769394a75",
		BotId:      CHATBOG.Id,
		CustomerId: CUSTOMERS[2].Id,
	},
}
