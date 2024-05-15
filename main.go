package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func PanicLogging(err error) {
	if err != nil {
		log.Panicln(err.Error())
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	PanicLogging(err)

	interval, err := strconv.Atoi(os.Getenv("MONITOR_INTERVAL_SEC"))
	PanicLogging(err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	PanicLogging(err)
	defer cli.Close()

	fmt.Println(" [*] Medic is on standby and transmitting")
	fmt.Printf(" [*] Interval Set: %d Second(s)\n", interval)

	for {
		containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
		PanicLogging(err)
		for _, container := range containers {
			cName := container.Names[0]
			cName = strings.ReplaceAll(cName, "/", "")
			cStatus := container.Status

			if strings.Contains(cStatus, "unhealthy") {
				fmt.Printf("%s - %s\n", cName, container.Status)
			}
		}
		time.Sleep(time.Second * time.Duration(interval))
	}

}
