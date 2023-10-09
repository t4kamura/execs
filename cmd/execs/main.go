package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/t4kamura/execs/internal/aws"
	"github.com/t4kamura/execs/internal/interactive"
)

const version = "0.0.0"

// TODO: panic -> error
func main() {
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		fmt.Printf("execs version %s\n", version)
		os.Exit(0)
	}

	profiles, err := aws.GetProfiles()
	if err != nil {
		panic(err)
	}
	if len(profiles) == 0 {
		panic("No profiles found")
	}

	profile, err := interactive.AskChoices("Select AWS profile", profiles)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cfg, err := aws.LoadConfig(ctx, profile)
	if err != nil {
		panic(err)
	}
	clnt := aws.New(&cfg)

	// cluster selection
	clusters, err := clnt.ListClusterNames(ctx)
	if err != nil {
		panic(err)
	}
	if len(clusters) == 0 {
		panic("No clusters found")
	}

	var selectedCluster string
	if len(clusters) == 1 {
		selectedCluster = clusters[0]
	} else {
		selectedCluster, err = interactive.AskChoices("Select cluster", clusters)
		if err != nil {
			panic(err)
		}
	}

	// service selection
	services, err := clnt.ListServiceNames(ctx, selectedCluster)
	if err != nil {
		panic(err)
	}
	if len(services) == 0 {
		panic("No services found")
	}

	var selectedService string
	if len(services) == 1 {
		selectedService = services[0]
	} else {
		selectedService, err = interactive.AskChoices("Select service", services)
		if err != nil {
			panic(err)
		}
	}

	// task selection
	tasks, err := clnt.ListTaskNames(ctx, selectedCluster, selectedService)
	if err != nil {
		panic(err)
	}
	if len(tasks) == 0 {
		panic("No services found")
	}

	var selectedTask string
	if len(tasks) == 1 {
		selectedTask = tasks[0]
	} else {
		selectedTask, err = interactive.AskChoices("Select task", tasks)
		if err != nil {
			panic(err)
		}
	}

	// container seleection
	containers, err := clnt.ListContainerNames(ctx, selectedCluster, selectedTask)
	if err != nil {
		panic(err)
	}
	if len(containers) == 0 {
		panic("No containers found")
	}

	var selectedContainer string
	if len(containers) == 1 {
		selectedContainer = containers[0]
	} else {
		selectedContainer, err = interactive.AskChoices("Select container", containers)
		if err != nil {
			panic(err)
		}
	}

	// start session
	if err := clnt.StartSession(ctx, selectedCluster, selectedTask, selectedContainer); err != nil {
		panic(err)
	}
}