package main

import (
	"encoding/json"
	"fmt"
	"github.com/ronna-s/scheduler/channels"
	. "github.com/ronna-s/scheduler/job"
	. "github.com/ronna-s/scheduler/scheduler"
	. "github.com/ronna-s/scheduler/workers"
	"time"
)

type (
	namedWorker struct {
		id int
		Worker
	}
)

func main() {
	go startScheduler()
	conf := &channels.ConsumerChannelConfig{
		ChannelConfig: channels.ChannelConfig{
			Name:     "jobs",
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		PrefetchCount: 1,
	}
	w1 := namedWorker{1, NewWorker(conf)}
	w2 := namedWorker{2, NewWorker(conf)}
	w3 := namedWorker{3, NewWorker(conf)}
	testFunc := func(nworker namedWorker) func(b []byte) error {
		return func(b []byte) error {
			fmt.Println(fmt.Sprintf("Worker id: %d: \t%s", nworker.id, string(b)))
			return nil
		}
	}

	go w1.HandleJobs(testFunc(w1))
	go w2.HandleJobs(testFunc(w2))
	go w3.HandleJobs(testFunc(w3))

	go publishJobsToSchduler()
	time.Sleep(25 * time.Second)

}
func startScheduler() {
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
func publishJobsToSchduler() {
	publisherCh := channels.NewPublisherChannel(&channels.PublisherChannelConfig{
		ChannelConfig: channels.ChannelConfig{
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		Exchange: "incoming",
	})
	body, _ := json.Marshal(Job{Data: []byte("1st - immediate"), Start: time.Now()})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("2nd - immediate - order can be confused with first"), Start: time.Now()})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("4th - 10 seconds from start"), Start: time.Now().Add(10 * time.Second)})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("last - 20 seconds from start"), Start: time.Now().Add(20 * time.Second)})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("3rd - 5 seconds from start"), Start: time.Now().Add(5 * time.Second)})
	publisherCh.Publish(body)

}
