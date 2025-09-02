// internal/types/struct.go

package types

type Config struct {
	Name                  string
	ReadTickerTimeSec     float32
	StatsPort             uint16
	MasterNode            bool
	MasterSend            bool
	MasterIP              string
	MasterPort            uint16
	MasterKey             string
	MasterSendIntervalMin float32
}
