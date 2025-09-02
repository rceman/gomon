package main

import "testing"

func TestLoadConfigDecimalTicker(t *testing.T) {
	t.Setenv("READ_TICKER_TIME_SEC", "0.5")
	t.Setenv("STATS_PORT", "9090")
	t.Setenv("VM_NAME", "TEST")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ReadTickerTimeSec != 0.5 {
		t.Fatalf("expected 0.5, got %v", cfg.ReadTickerTimeSec)
	}
}

func TestLoadConfigMasterSend(t *testing.T) {
	t.Setenv("READ_TICKER_TIME_SEC", "1")
	t.Setenv("STATS_PORT", "9090")
	t.Setenv("VM_NAME", "RNG")
	t.Setenv("MASTER_SEND", "true")
	t.Setenv("MASTER_SEND_INTERVAL_MIN", "1")
	t.Setenv("MASTER_IP", "127.0.0.1")
	t.Setenv("MASTER_PORT", "8080")
	t.Setenv("MASTER_KEY", "secret")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.MasterSend || cfg.MasterIP != "127.0.0.1" || cfg.MasterPort != 8080 || cfg.MasterKey != "secret" || cfg.MasterSendIntervalMin != 1 {
		t.Fatalf("config fields not loaded correctly: %+v", cfg)
	}
}
