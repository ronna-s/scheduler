package main

import (
	"encoding/json"
	. "github.com/ronna-s/scheduler/scheduler"
	"io/ioutil"
	"path"
	"runtime"
)

type ()

func main() {
	var config SchedulerConfig
	_, currentFilename, _, _ := runtime.Caller(0)
	cdir := path.Dir(currentFilename)
	file, err := ioutil.ReadFile(path.Join(cdir, "config.json"))
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	NewScheduler(config).Run()
}
