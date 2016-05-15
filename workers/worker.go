package workers

import (
	"github.com/ronna-s/scheduler/channels"
)

type (
	worker struct {
		conf channels.ConsumerChannelConfig
	}
	Worker interface {
		HandleJobs(cb func([]byte) error)
	}
)

func NewWorker(conf channels.ConsumerChannelConfig) Worker {
	return &worker{conf}
}

func (w *worker) HandleJobs(cb func([]byte) error) {
	incomingMessages := make(chan channels.Message)
	channels.NewConsumerAMQPChannel(w.conf).Listen(incomingMessages)
	for m := range incomingMessages {
		err := cb(m.Body())
		if err != nil {
			m.Reject()
		} else {
			m.Ack()
		}
	}
}
