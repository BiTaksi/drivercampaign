package eventmanager

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/BiTaksi/drivercampaign/pkg/event"
	"github.com/BiTaksi/drivercampaign/pkg/kafka"
)

type IEventManager interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) (*kafka.MessageInput, error)
	HandleException(eventType string, err error) error
}

type eventManager struct {
	handlerFactory IEventHandlerFactory
	eventFactory   event.IEventFactory
}

func NewEventManager(hFactory IEventHandlerFactory, tFactory event.IEventFactory) IEventManager {
	return &eventManager{
		handlerFactory: hFactory,
		eventFactory:   tFactory,
	}
}

func (em *eventManager) Handle(ctx context.Context, msg *sarama.ConsumerMessage) (*kafka.MessageInput, error) {
	attr := kafka.GetEventAttribute(msg.Headers)

	e, err := em.eventFactory.Make(msg.Topic, attr.Type, msg.Value)
	if err != nil {
		return nil, em.HandleException(attr.Type, err)
	}

	eh, ehErr := em.handlerFactory.Make(e, msg.Topic)
	if ehErr != nil {
		return nil, em.HandleException(attr.Type, ehErr)
	}

	input := &kafka.MessageInput{
		Topic:     msg.Topic,
		Key:       string(msg.Key),
		Value:     string(msg.Value),
		Attribute: attr,
	}

	if handleErr := eh.Handle(ctx); handleErr != nil {
		return input, em.HandleException(attr.Type, handleErr)
	}

	return input, nil
}

func (em *eventManager) HandleException(eventType string, err error) error {
	if err == event.ErrEventUnexpectedType {
		return nil
	}

	return event.ErrorBag{
		EventType:  eventType,
		Cause:      err,
		Reportable: true,
		Retryable:  em.isRetryable(err),
	}
}

func (em *eventManager) isRetryable(err error) bool {
	return err != event.ErrEventUnexpectedType && err != event.ErrHandlerUnexpectedType
}
