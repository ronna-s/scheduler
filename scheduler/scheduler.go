package scheduler

import (
	"encoding/json"
	"errors"
	"github.com/ronna-s/scheduler/channels"
	. "github.com/ronna-s/scheduler/job"
	"time"
)

type (
	Scheduler interface {
		Run()
	}
	SchedulerConfig struct {
		Publisher channels.PublisherChannelConfig
		Consumer  channels.ConsumerChannelConfig
	}

	scheduler struct {
		jobs          []Job
		incomingJobs  <-chan Job
		outGoing      chan Job
		consumerConf  channels.ConsumerChannelConfig
		publisherConf channels.PublisherChannelConfig
	}
)

func NewScheduler(conf SchedulerConfig) Scheduler {
	return &scheduler{consumerConf: conf.Consumer, publisherConf: conf.Publisher}
}
func (s *scheduler) Run() {
	incomingMessages := make(chan channels.Message)
	incomingJobs := make(chan Job, 1)
	outgoingJobs := make(chan Job, 1)
	channels.NewConsumerAMQPChannel(s.consumerConf).Listen(incomingMessages)
	publisherCh := channels.NewPublisherChannel(s.publisherConf)

	for {
		select {
		case message := <-incomingMessages:
			var job Job
			if err := json.Unmarshal(message.Body(), &job); err == nil {
				incomingJobs <- job
				message.Ack()
			} else {
				message.Reject()
			}
		case job := <-incomingJobs:
			s.pushJob(job)
		case <-time.After(s.timeUntilNextJob()):
			if job, err := s.popNextJob(); err == nil {
				outgoingJobs <- job
			} else {
				//log 1000 hours with no job
			}
		case job := <-outgoingJobs:
			publisherCh.Publish(job.Data)
		}
	}
}

func (s *scheduler) timeUntilNextJob() time.Duration {
	ts := 1000 * time.Hour
	for _, job := range s.jobs {
		if job.Start.Sub(time.Now()) < ts {
			ts = job.Start.Sub(time.Now())
		}
	}
	if ts < 0 {
		ts = 0
	}
	return ts
}
func (s *scheduler) pushJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *scheduler) popNextJob() (currJob Job, err error) {
	if len(s.jobs) == 0 {
		err = errors.New("no jobs in queue")
		return
	}
	idx := 0
	minTime := s.jobs[0].Start
	for i, job := range s.jobs {
		if job.Start.Unix() < minTime.Unix() {
			minTime = job.Start
			idx = i
		}
	}
	//pop item
	currJob = s.jobs[idx]
	s.jobs = append(s.jobs[0:idx], s.jobs[idx+1:]...)
	return
}
