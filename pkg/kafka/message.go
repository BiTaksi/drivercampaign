package kafka

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"

	constants "github.com/BiTaksi/drivercampaign/pkg/constans"

	"github.com/BiTaksi/drivercampaign/pkg/event"
)

type MessageInput struct {
	Topic     string
	Key       string
	Value     string
	Attribute *MessageAttribute
}

type MessageAttribute struct {
	Type      string
	Attempts  int
	CreatedAt int64
}

func NewMessageInput(topic, key string, e event.Event) *MessageInput {
	return &MessageInput{
		Topic:     topic,
		Key:       key,
		Value:     PreparePayload(e),
		Attribute: NewMessageAttribute(e.Type()),
	}
}

func NewMessageAttribute(t string) *MessageAttribute {
	return &MessageAttribute{
		Type:      t,
		Attempts:  0,
		CreatedAt: time.Now().Unix(),
	}
}

func GetEventAttribute(headers []*sarama.RecordHeader) *MessageAttribute {
	attr := new(MessageAttribute)
	for _, header := range headers {
		switch string(header.Key) {
		case constants.AttributeEventKey:
			attr.Type = string(header.Value)
		case constants.AttributeAttemptsKey:
			attr.Attempts = getAttempts(string(header.Value))
		case constants.AttributeCreatedAtKey:
			createdAt, _ := strconv.ParseInt(string(header.Value), constants.IntSmallBase, constants.Int64BitSize)
			attr.CreatedAt = createdAt
		}
	}

	return attr
}

func PreparePayload(e event.Event) string {
	j, _ := json.Marshal(e)
	return string(j)
}

func getAttempts(val string) int {
	attempts, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return attempts
}
