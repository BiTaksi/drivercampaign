package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type IConsumerInstance interface {
	Handler() ConsumerGroupHandler
	Consume(ctx context.Context, topics []string)
}

type consumerInstance struct {
	logger  *logrus.Logger
	client  sarama.ConsumerGroup
	handler ConsumerGroupHandler
}

func NewConsumerInstance(l *logrus.Logger, c sarama.ConsumerGroup, h ConsumerGroupHandler) IConsumerInstance {
	return &consumerInstance{
		logger:  l,
		client:  c,
		handler: h,
	}
}

func (k *consumerInstance) Handler() ConsumerGroupHandler {
	return k.handler
}

func (k *consumerInstance) Consume(ctx context.Context, topics []string) {
	for {
		if err := k.client.Consume(ctx, topics, k.handler); err != nil {
			k.logger.Fatalf("consume: %v", err)
		}

		if ctx.Err() != nil {
			return
		}

		k.handler.Ready()
	}
}
