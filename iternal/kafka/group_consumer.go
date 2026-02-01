package kafka

import (
	"context"
	"fmt"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
	"github.com/IBM/sarama"
)

type GroupConsumer struct {
	l       log.Logger
	groupID string
	addrs   []string
	config  *sarama.Config
	topics  []string
}

func NewGroupConsumer(groupID string, addrs []string, topics []string, config *sarama.Config, logger log.Logger) *GroupConsumer {
	return &GroupConsumer{
		groupID: groupID,
		addrs:   addrs,
		config:  config,
		topics:  topics,
	}
}

func (o *GroupConsumer) Run(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
	group, err := sarama.NewConsumerGroup(o.addrs, o.groupID, o.config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer func() {
		if err := group.Close(); err != nil {
			o.l.Error("failed to close consumer group", log.NewField("err", err))
		}
	}()

	go func() {
		for {
			select {
			case err, ok := <-group.Errors():
				if !ok {
					return
				}
				o.l.Error("kafka error", log.NewField("err", err))
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		err := group.Consume(ctx, o.topics, handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}
