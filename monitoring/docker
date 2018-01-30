package main

/*
	Check several docker elements such as:
	 	- node
		- container
*/
// NOTE: https://godoc.org/github.com/moby/moby/client

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Monitoring struct {
	message []string
	status  int
}

type SwarmStatus struct {
	ready   int
	manager int
	worker  int
	total   int
}

var (
	// Define context manager
	ctx = context.Background()

	// Log level
	INFO     = 0
	WARN     = 1
	CRITICAL = 2
)

func checkAndRaise(e error) {
	/*
		Raise exception when error is detected.
	*/
	if e != nil {
		panic(e)
	}

}

func log4Go(lvl int, message []string) {
	/*
		Print output with log level
	*/
	switch lvl {
	case WARN:
		fmt.Println("WARNNING -", strings.Join(message, ", "))
	case CRITICAL:
		fmt.Println("CRITICAL -", strings.Join(message, ", "))
	default:
		fmt.Println("INFO -", strings.Join(message, ", "))
	}

}

func main() {

	// Create arguments parser
	dockerEndpoint := flag.String("hosts", "unix:///var/run/docker.sock", "Docker endpoint.")
	dockerClientVersion := flag.String("client_version", "v1.24", "Docker client version.")
	swarmStatus := flag.Bool("swarm", false, "Get swarm status.")
	containerStatus := flag.Bool("containers", false, "Get containers status.")

	flag.Parse()

	// Create docker cli driver
	cli, err := client.NewClient(*dockerEndpoint, *dockerClientVersion, nil, nil)
	checkAndRaise(err)

	/*
		Get docker nodes statu.
	*/
	if *swarmStatus {
		// Create monitoring Object.
		monitoring := Monitoring{status: INFO}

		// Define node ready counter.
		swarmStatus := SwarmStatus{0, 0, 0, 0}

		// get Node list
		nodes, err := cli.NodeList(ctx, types.NodeListOptions{})
		checkAndRaise(err)

		// Check node elements.
		for _, node := range nodes {
			// Increment total
			swarmStatus.total++

			// Count node with ready status
			if node.Status.State == "ready" {
				swarmStatus.ready++
			}

			// Count manager and worker
			if node.Spec.Role == "manager" {
				swarmStatus.manager++
			} else if node.Spec.Role == "worker" {
				swarmStatus.worker++
			}

		}

		// Set monitoring Status
		monitoring.message = []string{fmt.Sprintf("%d/%d nodes ready", swarmStatus.ready, swarmStatus.total),
			fmt.Sprintf("%d/%d manager", swarmStatus.manager, swarmStatus.total),
			fmt.Sprintf("%d/%d worker", swarmStatus.manager, swarmStatus.total)}

		log4Go(monitoring.status, monitoring.message)
	}

	if *containerStatus {

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		checkAndRaise(err)

		for _, container := range containers {
			fmt.Printf("%s %s\n", container.Names, container.Status)

			/*
				containerStats, err := cli.ContainerStats(context.Background(), container.ID[:10], true)
				checkAndRaise(err)
				buf := new(bytes.Buffer)
				buf.ReadFrom(containerStats.Body)
				strStats := buf.String()
				fmt.Printf("%s\n", strStats)
			*/
		}
	}
}
