package main

import (
	"fmt"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda"
	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"math/rand"
	"os"
	"time"
)

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	command := app.NewSchedulerCommand(
		app.WithPlugin(yoda.Name, yoda.New),
	)
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}