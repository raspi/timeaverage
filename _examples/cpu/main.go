package main

import (
	"bufio"
	"fmt"
	"github.com/raspi/timeaverage"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getCPUUsage() (usage float64, err error) {
	f, err := os.Open(`/proc/stat`)

	if err != nil {
		return 0.0, err
	}

	scanner := bufio.NewScanner(f)

	var cpu []float64

	for scanner.Scan() {
		txt := scanner.Text()

		if !strings.HasPrefix(txt, `cpu `) {
			continue
		}

		for idx, v := range strings.Split(txt, ` `) {
			if idx == 0 || idx == 1 {
				continue
			}

			val, err := strconv.Atoi(v)
			if err != nil {
				return 0.0, err
			}

			cpu = append(cpu, float64(val))

		}

		break

	}

	l := len(cpu)

	if l != 10 {
		return 0.0, fmt.Errorf(`invalid amount of fields: %d`, l)
	}

	usage = (cpu[0] + cpu[2]) * 100.0 / (cpu[0] + cpu[2] + cpu[3])

	return usage, nil
}

func main() {
	ta := timeaverage.New(time.Second*10, time.Millisecond*500, 0.0, getCPUUsage)

	ta.Start()

	for {
		usage := ta.Average()
		log.Printf(`cpu usage: %f`, usage)
		time.Sleep(time.Second * 1)
	}

}
