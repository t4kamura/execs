package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/t4kamura/execs/internal/aws"
	"github.com/t4kamura/execs/internal/interactive"
)

const version = "0.0.1"

func main() {
	v := flag.Bool("v", false, "show version")
	s := flag.String("s", "/bin/bash", "shell used within ECS container")
	flag.Parse()

	if *v {
		fmt.Printf("execs version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(1)
	}

	// profile selection
	profiles, err := aws.GetProfiles()
	if err != nil {
		log.Fatal(err)
	}
	if len(profiles) == 0 {
		log.Fatal("No AWS profiles found")
	}

	profile, err := interactive.SelectItem("Select AWS profile", profiles)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AWS profile: %s\n", profile)

	ctx := context.Background()
	cfg, err := aws.LoadConfig(ctx, profile)
	if err != nil {
		log.Fatal(err)
	}
	clnt := aws.New(&cfg)

	// cluster selection
	clusters, err := clnt.ListClusterNames(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(clusters) == 0 {
		log.Fatal("No clusters found")
	}

	var selectedCluster string
	if len(clusters) == 1 {
		selectedCluster = clusters[0]
		fmt.Printf("Cluster: %s (auto selected)\n", selectedCluster)
	} else {
		selectedCluster, err = interactive.SelectItem("Select cluster", clusters)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Cluster: %s\n", selectedCluster)
	}

	// service selection
	services, err := clnt.ListServiceNames(ctx, selectedCluster)
	if err != nil {
		log.Fatal(err)
	}
	if len(services) == 0 {
		log.Fatal("No services found")
	}

	var selectedService string
	if len(services) == 1 {
		selectedService = services[0]
		fmt.Printf("Service: %s (auto selected)\n", selectedService)
	} else {
		selectedService, err = interactive.SelectItem("Select service", services)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Service: %s\n", selectedService)
	}

	// task selection
	tasks, err := clnt.ListTaskNames(ctx, selectedCluster, selectedService)
	if err != nil {
		log.Fatal(err)
	}
	if len(tasks) == 0 {
		log.Fatal("No services found")
	}

	var selectedTask string
	if len(tasks) == 1 {
		selectedTask = tasks[0]
		fmt.Printf("Task: %s (auto selected)\n", selectedTask)
	} else {
		selectedTask, err = interactive.SelectItem("Select task", tasks)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task: %s\n", selectedTask)
	}

	// container seleection
	containers, err := clnt.ListContainerNames(ctx, selectedCluster, selectedTask)
	if err != nil {
		log.Fatal(err)
	}
	if len(containers) == 0 {
		log.Fatal("No containers found")
	}

	var selectedContainer string
	if len(containers) == 1 {
		selectedContainer = containers[0]
		fmt.Printf("Container: %s (auto selected)\n", selectedContainer)
	} else {
		selectedContainer, err = interactive.SelectItem("Select container", containers)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Container: %s\n", selectedContainer)
	}

	// start session
	if err := clnt.StartSession(ctx, selectedCluster, selectedTask, selectedContainer, *s, cfg.Region); err != nil {
		log.Fatal(err)
	}
}
