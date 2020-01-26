package main

import (
	"fmt"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/register"
	"k8s.io/component-base/logs"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	command := register.Register()
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
