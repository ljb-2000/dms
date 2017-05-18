package cmd

import (
	"fmt"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	"time"
)

// ContainersLogsCmd show container logs
func ContainersLogsCmd(id string) {
	logs, err := GetContainersLogs(id)
	if err != nil {
		fmt.Println("Oops! Something went wrong, please try again")
		return
	}

	fmt.Println(logs)
}

// ContainersMetricsCmd show container(s) metrics
func ContainersMetricsCmd(id []string) {
	metrics, err := GetContainersMetrics(strings.Join(id, " "))
	if err != nil {
		fmt.Println("Oops! Something went wrong, please try again")
		return
	}

	// if metrics length == 0 -> no running containers
	if len(metrics) == 0 {
		fmt.Println("No running containers")
		return
	}

	// print metrics table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CPU, %", "MEM, %", "TIME"})
	for _, m := range metrics {
		table.Append([]string{
			m.Name,
			strconv.FormatFloat(m.CPUPercentage, 'f', 2, 64),
			strconv.FormatFloat(m.MemoryPercentage, 'f', 2, 64),
			time.Now().String()[0:19],
		})
	}
	table.Render()
}

// StoppedContainersCmd show stopped containers
func StoppedContainersCmd() {
	stopped, err := GetStoppedContainers()
	if err != nil {
		fmt.Println("Oops! Something went wrong, please try again")
		return
	}

	// if first stopped array element == "no stopped containers"
	// -> no stopped containers
	if stopped[0] == "No stopped containers" {
		fmt.Println(stopped[0])
		return
	}

	i := 0
	for _, s := range stopped {
		fmt.Println(s)
		i++
	}
	fmt.Println("Total stopped:", i)
}

// LaunchedContainersCmd show launched containers
func LaunchedContainersCmd() {
	launched, err := GetLaunchedContainers()
	if err != nil {
		fmt.Println("Oops! Something went wrong, please try again")
		return
	}

	// if first launched array element == "no launched containers"
	// -> no launched containers
	if launched[0] == "No launched containers" {
		fmt.Println(launched[0])
		return
	}

	i := 0
	for _, l := range launched {
		fmt.Println(l)
		i++
	}
	fmt.Println("Total launched:", i)
}
