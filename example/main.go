package main

import (
	"encoding/json"
	"fmt"
	"github.com/ronna-s/scheduler/channels"
	. "github.com/ronna-s/scheduler/job"
	. "github.com/ronna-s/scheduler/scheduler"
	. "github.com/ronna-s/scheduler/workers"
	"io/ioutil"
	"path"
	"runtime"
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

	conf := channels.ConsumerChannelConfig{
		ChannelConfig: channels.ChannelConfig{
			Name:     "jobs",
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		PrefetchCount: 1,
	}
	fmt.Println("Starting 3 workers\n>>>>>>>>>>")
	w1 := namedWorker{1, NewWorker(conf)}
	w2 := namedWorker{2, NewWorker(conf)}
	w3 := namedWorker{3, NewWorker(conf)}
	st := time.Now()
	testFunc := func(nworker namedWorker) func(b []byte) error {
		return func(b []byte) error {
			fmt.Println(
				fmt.Sprintf(
					"Worker #%d received job. Time since beginning: %v\nJob info: \n\t%s\n====================",
					nworker.id, time.Since(st), string(b)))
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
	var config SchedulerConfig
	_, currentFilename, _, _ := runtime.Caller(0)
	cdir := path.Dir(currentFilename)
	file, err := ioutil.ReadFile(path.Join(cdir, "../config.json"))
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	NewScheduler(config).Run()
}
func publishJobsToSchduler() {
	publisherCh := channels.NewPublisherChannel(channels.PublisherChannelConfig{
		ChannelConfig: channels.ChannelConfig{
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		Exchange: "incoming",
	})
	fmt.Println("Publishing 5 jobs for 20 seconds\n<<<<<<<<<\n")
	body, _ := json.Marshal(Job{Data: []byte("execute me immediatately"), Start: time.Now()})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("execute me immediatately"), Start: time.Now()})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("execute me after 10 seconds"), Start: time.Now().Add(10 * time.Second)})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("execute me after 20 seconds"), Start: time.Now().Add(20 * time.Second)})
	publisherCh.Publish(body)
	body, _ = json.Marshal(Job{Data: []byte("execute me after 5 seconds"), Start: time.Now().Add(5 * time.Second)})
	publisherCh.Publish(body)

}
