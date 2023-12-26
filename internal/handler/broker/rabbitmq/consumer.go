package rabbitmq

import (
	"encoding/json"
	"github.com/Verce11o/yata-notifications/internal/domain"
	"github.com/Verce11o/yata-notifications/internal/handler/websockets"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NotificationConsumer struct {
	AmqpConn *amqp.Connection
	log      *zap.SugaredLogger
	trace    trace.Tracer
}

func NewNotificationConsumer(amqpConn *amqp.Connection, log *zap.SugaredLogger, trace trace.Tracer) *NotificationConsumer {
	return &NotificationConsumer{AmqpConn: amqpConn, log: log, trace: trace}
}

func (c *NotificationConsumer) createChannel(exchangeName, queueName, bindingKey string) *amqp.Channel {
	ch, err := c.AmqpConn.Channel()

	if err != nil {
		panic(err)
	}

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

func (c *NotificationConsumer) StartConsumer(queueName, consumerTag, exchangeName, bindingKey string, clients websockets.WsClients) (<-chan amqp.Delivery, error) {
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
		go c.worker(i, deliveries, clients)
	}
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.log.Infof("Notify close: %v", chanErr)

	return deliveries, chanErr

}

func (c *NotificationConsumer) worker(index int, messages <-chan amqp.Delivery, clients websockets.WsClients) {
	for message := range messages {
		c.log.Infof("Worker #%d: %v", index, string(message.Body))

		var request domain.IncomingNewTweetNotification

		err := json.Unmarshal(message.Body, &request)

		if err != nil {
			c.log.Errorf("failed to unmarshal request: %v", err)
		}

		err = message.Ack(false)

		if err != nil {
			c.log.Errorf("failed to acknowledge delivery: %v", err)
		}

	}
	c.log.Info("Channel closed")
}
