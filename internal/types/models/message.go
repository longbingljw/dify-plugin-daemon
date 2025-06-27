package modelsAdd commentMore actions

import "time"

type Message struct {
	Channel   string    `json:"channel" gorm:"column:channel;type:varchar(1024);index: idx_channel_id,unique;not null"`
	ID        int64     `json:"id" gorm:"column:id;primary_key;type:bigint(20) auto_increment;not null;index: idx_channel_id;comment:'ID';"`
	Message   string    `json:"message" gorm:"column:message;type:varchar(1024);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageSubscribe struct {
	Channel       string    `json:"channel" gorm:"column:channel;type:varchar(1024);uniqueIndex:idx_channel_subscriber;not null"`
	Subscriber    string    `json:"subscriber" gorm:"column:subscriber;type:varchar(1024);uniqueIndex:idx_channel_subscriber;not null"`
	LastMessageId int64     `json:"lastMessageId" gorm:"column:lastMessageId;type:bigint(20);not null;comment:'lastMessageId';default:-1"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}