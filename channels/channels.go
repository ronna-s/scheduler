package channels

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type (
	ChannelConfig struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
	ConsumerChannelConfig struct {
		ChannelConfig
		PrefetchCount int
	}
	PublisherChannelConfig struct {
		ChannelConfig
		Exchange         string
		RoutingKeyPrefix string
	}
	Channel interface {
		Listen(messages chan<- Message)
	}
	PublisherChannel interface {
		Publish([]byte)
	}
	aMQPChannel struct {
		channel *amqp.Channel
	}
	aMQPPublisherChannel struct {
		aMQPChannel
		Exchange         string
		RoutingKeyPrefix string
	}
	AMQPChannel interface {
		Channel
	}
	aMQPConsumerChannel struct {
		aMQPChannel
		prefetchCount int
		queue         string
	}

	Delivery interface {
		Ack(bool) error
		Reject(bool) error
	}
	ackableMessage struct {
		body     []byte
		delivery Delivery
	}
	Message struct {
		ackableMessage
	}
)

func getAmqpChannel(url string) (*amqp.Channel, error) {
	var (
		c   *amqp.Connection
		ch  *amqp.Channel
		err error
	)

	if c, err = amqp.Dial(url); err == nil {
		if ch, err = c.Channel(); err != nil {
			c.Close()
		}
	} else {
		panic(fmt.Sprintf("[AMPQ] attempting to re-connect to amqp server. Error: %s", err))
	}

	return ch, err
}
func NewConsumerAMQPChannel(conf *ConsumerChannelConfig) AMQPChannel {
	url := "amqp://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port
	channel, err := getAmqpChannel(url)
	if err != nil {
		panic(err)
	}
	return &aMQPConsumerChannel{
		aMQPChannel:   aMQPChannel{channel},
		prefetchCount: conf.PrefetchCount,
		queue:         conf.Name,
	}

}
func (ch *aMQPConsumerChannel) Listen(messages chan<- Message) {
	ch.channel.Qos(ch.prefetchCount, 0, false)
	deliveries, err := ch.channel.Consume(ch.queue, "", false, false, false, false, nil)
	if err == nil {
		go func(deliveries <-chan amqp.Delivery, ch AMQPChannel) {
			for delivery := range deliveries {
				messages <- Message{
					ackableMessage{
						body:     delivery.Body,
						delivery: delivery,
					},
				}
			}
		}(deliveries, ch)
	} else {
		panic("failed listening to amqp channel")
	}
	return
}
func (m *ackableMessage) Body() []byte {
	return m.body
}
func (m *ackableMessage) Ack() error {
	return m.delivery.Ack(false)
}
func (m *ackableMessage) Reject() error {
	return m.delivery.Reject(false)
}

func NewPublisherChannel(conf *PublisherChannelConfig) PublisherChannel {
	url := "amqp://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port
	channel, err := getAmqpChannel(url)
	if err != nil {
		panic(err)
	}
	return &aMQPPublisherChannel{
		aMQPChannel:      aMQPChannel{channel},
		Exchange:         conf.Exchange,
		RoutingKeyPrefix: conf.RoutingKeyPrefix,
	}

}
func (p *aMQPPublisherChannel) Publish(body []byte) {
	err := p.channel.Publish(p.Exchange, p.RoutingKeyPrefix, true, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		panic(err)
	}
}
