package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/Verce11o/yata-notifications/internal/domain"
	"github.com/Verce11o/yata-notifications/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NotificationConsumer struct {
	AmqpConn *amqp.Connection
	log      *zap.SugaredLogger
	tracer   trace.Tracer
	service  service.Notifications
}

func NewNotificationConsumer(amqpConn *amqp.Connection, log *zap.SugaredLogger, trace trace.Tracer, service service.Notifications) *NotificationConsumer {
	return &NotificationConsumer{AmqpConn: amqpConn, log: log, tracer: trace, service: service}
}

func (c *NotificationConsumer) createChannel(exchangeName, queueName, bindingKey string) *amqp.Channel {
	ch, err := c.AmqpConn.Channel()

	if err != nil {
		panic(err)
	}

	// think about changing its kind
	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return ch

}

func (c *NotificationConsumer) StartConsumer(queueName, consumerTag, exchangeName, bindingKey string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := c.createChannel(exchangeName, queueName, bindingKey)
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		i := i
		go c.worker(ctx, i, deliveries)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.log.Infof("Notify close: %v", chanErr)

	return chanErr

}

func (c *NotificationConsumer) worker(ctx context.Context, index int, messages <-chan amqp.Delivery) {
	for message := range messages {
		c.log.Infof("Worker #%d: %v", index, string(message.Body))

		var request domain.IncomingNewNotification

		err := json.Unmarshal(message.Body, &request)

		if err != nil {
			c.log.Errorf("failed to unmarshal request: %v", err)
			if err := message.Nack(false, false); err != nil {
				c.log.Errorf("cannot nack message: %v", err)
			}
			return
		}

		subscribers, err := c.service.GetUserSubscribers(ctx, request.SenderID.String())
		c.log.Debugf("%#v", subscribers)

		if err != nil {
			c.log.Errorf("failed to get user subscribers: %v", err)
			if err := message.Nack(false, false); err != nil {
				c.log.Errorf("cannot nack message: %v", err)
			}
			return
		}

		err = c.service.BatchAddNotification(ctx, subscribers, request)

		if err != nil {
			c.log.Errorf("failed to add notification: %v", err)
			if err := message.Nack(false, false); err != nil {
				c.log.Errorf("cannot nack message: %v", err)
			}
			return
		}

		err = message.Ack(false)

		if err != nil {
			c.log.Errorf("failed to acknowledge delivery: %v", err)
		}

	}
	c.log.Info("Channel closed")
}
