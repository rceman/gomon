package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/therceman/gomon/internal/app"
	"github.com/therceman/gomon/internal/dotenv"
	"github.com/therceman/gomon/internal/helpers"
	"github.com/therceman/gomon/internal/types"
)

// LoadConfig loads environment variables into a Config struct
func LoadConfig() (types.Config, error) {
	readTickerTimeSec, err := helpers.ConvertStringToUint16(os.Getenv("READ_TICKER_TIME_SEC"))
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid value for READ_TICKER_TIME_SEC")
	}

	flushTickerTimeSec, err := helpers.ConvertStringToUint16(os.Getenv("FLUSH_TICKER_TIME_SEC"))
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid value for FLUSH_TICKER_TIME_SEC")
	}

	sleepBetweenFetchesMs, err := helpers.ConvertStringToUint16(os.Getenv("SLEEP_BETWEEN_FETCHES_MS"))
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid value for SLEEP_BETWEEN_FETCHES_MS")
	}

	var metricKeys []string
	keys := os.Getenv("METRIC_KEYS")
	if keys == "" {
		metricKeys = []string{
			"id", "name", "cpu_max_perc", "cpu_avg_perc", "mem_max_mb", "mem_avg_mb", "mem_max_perc", "mem_avg_perc", "disk_mb",
		}
	} else {
		metricKeys = strings.Split(keys, ",")
	}

	config := types.Config{
		ContainerName:         os.Getenv("CONTAINER_NAME"),
		GrafanaInfluxURL:      os.Getenv("GRAFANA_INFLUX_URL"),
		GrafanaAPIKey:         os.Getenv("GRAFANA_API_KEY"),
		GrafanaUsername:       os.Getenv("GRAFANA_USERNAME"),
		ReadTickerTimeSec:     readTickerTimeSec,
		FlushTickerTimeSec:    flushTickerTimeSec,
		SleepBetweenFetchesMs: sleepBetweenFetchesMs,
		MetricKeys:            metricKeys,
	}

	return config, nil
}

func main() {
	// Load the environment variables from the .env file
	err := dotenv.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}

	// Load the configuration into the Config struct
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Pass the Config struct to the app
	app.Run(config)
}
