package queueAdd commentMore actions

import (
	"time"

	"github.com/langgenius/dify-plugin-daemon/internal/db"
	"github.com/langgenius/dify-plugin-daemon/internal/types/models"
	"github.com/langgenius/dify-plugin-daemon/internal/utils/log"
	"github.com/langgenius/dify-plugin-daemon/internal/utils/parser"
)

func init() {
	go cleanMessages()
}

func cleanMessages() {
	time.Sleep(time.Minute * 5)
	for true {
		log.Info("cleaning messages and message_subscribes")
		db.Exec("delete from messages where created_at < now() - interval ? hour;", 1)
		db.Exec("delete from message_subscribes where updated_at < now() - interval ? hour;", 1)
		time.Sleep(time.Minute * 1)
	}
}
func ReplaceInto(data models.MessageSubscribe) error {
	return db.Exec("replace into message_subscribes values (?)",
		data.Channel, data.Subscriber, data.LastMessageId, data.CreatedAt, time.Now())
}

func Subscribe[T any](channel string, clusterId string) (<-chan T, func()) {
	ch := make(chan T)
	_, err := db.GetOne[models.MessageSubscribe](
		db.Equal("channel", channel),
		db.Equal("subscriber", clusterId),
	)

	if err != nil {
		ms := &models.MessageSubscribe{
			Channel:       channel,
			Subscriber:    clusterId,
			LastMessageId: -1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		log.Info("failed to get the message subscribe channel: %s, create a new record", err.Error())
		ReplaceInto(*ms)
	}
	go func() {
		defer close(ch)
		alive := true
		for alive {
			messageSubscribe, _ := db.GetOne[models.MessageSubscribe](
				db.Equal("channel", channel),
				db.Equal("subscriber", clusterId),
			)
			lastMessageId := messageSubscribe.LastMessageId
			createTime := messageSubscribe.CreatedAt
			//log.Debug("last message id: %d", lastMessageId)
			message, err := db.GetAny[models.Message]("select * from messages where channel = ? and id > ? and created_at > ? order by id asc limit 1", channel, lastMessageId, createTime)
			//log.Debug("message:%s", message.ID)
			if err != nil {
				log.Warn("failed to receive message from ob, will retry in 1 second")
				time.Sleep(1 * time.Second)
				continue
			}
			if message.ID < 1 {
				time.Sleep(5 * time.Second)
			} else {
				v, err := parser.UnmarshalJson[T](message.Message)
				if err != nil {
					log.Error("failed to unmarshal message from ob: %s, will retry in 1 second", err.Error())
				} else {
					ch <- v
				}
			}

			newRecord := &models.MessageSubscribe{
				Channel:       channel,
				Subscriber:    clusterId,
				LastMessageId: message.ID,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			//log.Debug("new Record:%s", newRecord.LastMessageId)
			ReplaceInto(*newRecord)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return ch, func() {
	}
}

func Publish(channel string, msg any) error {
	jsonStr := ""
	if _, ok := msg.(string); !ok {
		jsonStr = parser.MarshalJson(msg)
	}
	message := &models.Message{
		Channel:   channel,
		Message:   jsonStr,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(message)
	return nil
}