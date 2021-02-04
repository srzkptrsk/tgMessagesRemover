package main

import (
	"github.com/jinzhu/gorm"
)

// Message structure
type Message struct {
	gorm.Model
	ChatId    int64
	MessageId int
}
