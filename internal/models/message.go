package models

import (
	"time"
)

type Message struct {
	Id          string    `json:"id"`
	Body        string    `json:"body"`
	ChatId      string    `json:"chatId"`
	From        string    `json:"from"`
	TimeCreated time.Time `json:"timeCreated"`
}
