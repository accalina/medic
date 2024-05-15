package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
				fmt.Printf("%s - %s\n", cName, cStatus)
				t := time.Now()
				msg := fmt.Sprintf("Container '%s' status is '%s', TIME: %s", cName, cStatus, t.Format("2006-01-02 15:04:05"))
				sendToNtfy(msg)
			}
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func sendToNtfy(message string) {
	payload := strings.NewReader(message)

	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("NTFY_URL"), payload)
	PanicLogging(err)

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Set("Title", os.Getenv("NTFY_TITLE"))
	req.Header.Set("Priority", os.Getenv("NTFY_PRIORITY"))
	req.Header.Set("Tags", os.Getenv("NTFY_TAGS"))
	res, err := client.Do(req)
	PanicLogging(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	PanicLogging(err)

	fmt.Println(string(body))
}
