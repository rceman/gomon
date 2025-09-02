package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/therceman/gomon/internal/app"
	"github.com/therceman/gomon/internal/dotenv"
	"github.com/therceman/gomon/internal/helpers"
	"github.com/therceman/gomon/internal/types"
)

func LoadConfig() (types.Config, error) {
	readTickerTimeSec, err := helpers.ConvertStringToFloat32(os.Getenv("READ_TICKER_TIME_SEC"))
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid value for READ_TICKER_TIME_SEC")
	}

	statsPort, err := helpers.ConvertStringToUint16(os.Getenv("STATS_PORT"))
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid value for STATS_PORT")
	}

	name := os.Getenv("VM_NAME")
	if name == "" {
		return types.Config{}, fmt.Errorf("VM_NAME is required")
	}

	masterNode, _ := strconv.ParseBool(os.Getenv("MASTER_NODE"))
	masterSend, _ := strconv.ParseBool(os.Getenv("MASTER_SEND"))

	masterPort, _ := helpers.ConvertStringToUint16(os.Getenv("MASTER_PORT"))

	masterSendIntervalMin := float32(0.5)
	if v := os.Getenv("MASTER_SEND_INTERVAL_MIN"); v != "" {
		if parsed, err := helpers.ConvertStringToFloat32(v); err == nil {
			masterSendIntervalMin = parsed
		}
	}

	config := types.Config{
		Name:                  name,
		ReadTickerTimeSec:     readTickerTimeSec,
		StatsPort:             statsPort,
		MasterNode:            masterNode,
		MasterSend:            masterSend,
		MasterIP:              os.Getenv("MASTER_IP"),
		MasterPort:            masterPort,
		MasterKey:             os.Getenv("MASTER_KEY"),
		MasterSendIntervalMin: masterSendIntervalMin,
	}

	return config, nil
}

func main() {
	err := dotenv.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	app.Run(config)
}
