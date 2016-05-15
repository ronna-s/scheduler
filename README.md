# scheduler solution based on rabbitmq
scheduler receives jobs from rabbit and publishes jobs upon start time to exchange based on config
workers receive jobs data from rabbit based on config and run them

####setup test environment
`vagrant up`

`vagrant ssh`

####inside vagrant - run example main
`go run $SCHEDULER_HOME_DIR/main.go`


