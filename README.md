# scheduler solution based on rabbitmq
scheduler receives jobs from rabbit and publishes jobs upon start time to exchange based on config
workers receive jobs data from rabbit based on config and run them

####setup test environment (vagrant)
`git clone https://github.com/ronna-s/scheduler && cd scheduler && vagrant up && vagrant ssh`

to run example main:
`go run $SCHEDULER_HOME_DIR/example/main.go`
