# scheduler solution based on rabbitmq
scheduler receives jobs from rabbit and publishes jobs upon start time to exchanges for workers (consumers) based on config
workers receive jobs data from rabbit based on config and run them

####setup test environment (vagrant)
`git clone https://github.com/ronna-s/scheduler && cd scheduler && vagrant up && vagrant ssh`

to run example main:
`go run $SCHEDULER_HOME_DIR/example/main.go`

###Setup scheduler by editing config.json

###Setup workers to handle jobs:
Hand over a callback to run to `worker.HandleJobs`

```go
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
  NewWorker(amqpChConfig).HandleJobs(func(b []byte) (err error) {
    var d myDataType
    if err = json.Unmarshal(b,&d); err!=nil{
      return
    }
    return foo(d)
  }

```
