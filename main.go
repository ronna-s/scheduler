package main

import (
	"github.com/ronna-s/scheduler/channels"
	. "github.com/ronna-s/scheduler/scheduler"
)

func main() {
	NewScheduler(&channels.ConsumerChannelConfig{
		ChannelConfig: channels.ChannelConfig{
			Name:     "incoming",
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		PrefetchCount: 1,
	},
		&channels.PublisherChannelConfig{
			ChannelConfig: channels.ChannelConfig{
				User:     "guest",
				Password: "guest",
				Host:     "localhost",
				Port:     "5672",
			},
			Exchange: "jobs",
		}).Run()
}
