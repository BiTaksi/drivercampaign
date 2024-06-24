package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type syncHandler struct {
	logger          *logrus.Logger
	consumerManager IConsumerManager

	ready chan bool
}

func NewSyncHandler(l *logrus.Logger, cm IConsumerManager) ConsumerGroupHandler {
	return &syncHandler{
		logger:          l,
		consumerManager: cm,
		ready:           make(chan bool),
	}
}

func (h *syncHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *syncHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *syncHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer h.logger.WithField("memberId", sess.MemberID()).Info("consumer: sync group session closed")

	for {
		select {
		case <-sess.Context().Done():
			return nil
		case msg := <-claim.Messages():
			if msg == nil || sess.Context().Err() != nil {
				return nil
			}

			if err := h.consumerManager.Process(sess.Context(), msg); err == nil {
				sess.MarkMessage(msg, "")
			}
		}
	}
}

func (h *syncHandler) Ready() {
	h.ready = make(chan bool)
}

func (h *syncHandler) Status() chan bool {
	return h.ready
}
