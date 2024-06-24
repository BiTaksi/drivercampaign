package kafka

import (
	"strconv"

	"github.com/IBM/sarama"

	constants "github.com/BiTaksi/drivercampaign/pkg/constans"
)

type Output struct {
	Partition int32
	Offset    int64
}

type IProducerInstance interface {
	SendMessage(mi *MessageInput) (*Output, error)
	PrepareMessage(mi *MessageInput) *sarama.ProducerMessage
	Close() error
}

type producerInstance struct {
	producer sarama.SyncProducer
}

func InitProducer(config *sarama.Config, brokers []string) (IProducerInstance, error) {
	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &producerInstance{producer: p}, nil
}

func (p *producerInstance) SendMessage(mi *MessageInput) (*Output, error) {
	partition, offset, err := p.producer.SendMessage(p.PrepareMessage(mi))
	if err != nil {
		return nil, err
	}

	return &Output{
		Partition: partition,
		Offset:    offset,
	}, nil
}

func (p *producerInstance) PrepareMessage(mi *MessageInput) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(mi.Key),
		Topic: mi.Topic,
		Value: sarama.StringEncoder(mi.Value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(constants.AttributeEventKey),
				Value: []byte(mi.Attribute.Type),
			},
			{
				Key:   []byte(constants.AttributeAttemptsKey),
				Value: []byte(strconv.Itoa(mi.Attribute.Attempts)),
			},
			{
				Key:   []byte(constants.AttributeCreatedAtKey),
				Value: []byte(strconv.FormatInt(mi.Attribute.CreatedAt, constants.IntSmallBase)),
			},
		},
	}
}

func (p *producerInstance) Close() error {
	return p.producer.Close()
}
