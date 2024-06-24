package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"

	"github.com/BiTaksi/drivercampaign/pkg/nrclient"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler

	Ready()
	Status() chan bool
}

type CustomHandler interface {
	Do(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type IConsumerManager interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type ConsumerSessionMessage struct {
	Session sarama.ConsumerGroupSession
	Message *sarama.ConsumerMessage
}

type consumerManager struct {
	logger           *logrus.Logger
	customHandler    CustomHandler
	newRelicInstance nrclient.INewRelicInstance
}

func NewConsumerManager(l *logrus.Logger, ch CustomHandler, ni nrclient.INewRelicInstance) IConsumerManager {
	return &consumerManager{
		logger:           l,
		customHandler:    ch,
		newRelicInstance: ni,
	}
}

func (cm *consumerManager) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	txn := cm.newRelicInstance.Application().StartTransaction(fmt.Sprintf("key:%s", string(msg.Key)))
	txn.AddAttribute("event.topic", msg.Topic)
	txn.AddAttribute("event.partition", msg.Partition)
	txn.AddAttribute("event.offset", msg.Offset)

	defer txn.End()

	ctx = newrelic.NewContext(ctx, txn)
	if err := cm.customHandler.Do(ctx, msg); err != nil {
		cm.handleException(msg, txn, err)
	}

	return nil
}

func (cm *consumerManager) prepareLogFields(msg *sarama.ConsumerMessage) logrus.Fields {
	return logrus.Fields{
		"topic":     msg.Topic,
		"key":       string(msg.Key),
		"partition": msg.Partition,
		"offset":    msg.Offset,
		"body":      string(msg.Value),
	}
}

func (cm *consumerManager) handleException(msg *sarama.ConsumerMessage, txn *newrelic.Transaction, err error) {
	fields := logrus.Fields{
		"event": cm.prepareLogFields(msg),
	}

	eb, ok := err.(ErrorBag)
	if !ok {
		txn.NoticeError(err)
		cm.logger.WithFields(fields).WithError(err).Errorf("event: key: %s - message: processing error", string(msg.Key))
		return
	}

	if !eb.IsReportable() {
		cm.logger.WithFields(fields).WithError(err).Warn(eb.FormattedText())
		return
	}

	cm.logger.WithFields(fields).WithError(err).Error(eb.FormattedText())

	txn.NoticeError(eb)
}
