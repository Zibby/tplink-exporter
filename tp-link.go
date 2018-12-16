package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sausheong/hs1xxplug"
	log "github.com/sirupsen/logrus"
)

type kasa struct {
	System struct {
		GetSysinfo struct {
			ErrCode    int     `json:"err_code"`
			SwVer      string  `json:"sw_ver"`
			HwVer      string  `json:"hw_ver"`
			Type       string  `json:"type"`
			Model      string  `json:"model"`
			Mac        string  `json:"mac"`
			DeviceID   string  `json:"deviceId"`
			HwID       string  `json:"hwId"`
			FwID       string  `json:"fwId"`
			OemID      string  `json:"oemId"`
			Alias      string  `json:"alias"`
			DevName    string  `json:"dev_name"`
			IconHash   string  `json:"icon_hash"`
			RelayState int     `json:"relay_state"`
			OnTime     int     `json:"on_time"`
			ActiveMode string  `json:"active_mode"`
			Feature    string  `json:"feature"`
			Updating   int     `json:"updating"`
			Rssi       int     `json:"rssi"`
			LedOff     int     `json:"led_off"`
			Latitude   float64 `json:"latitude"`
			Longitude  float64 `json:"longitude"`
		} `json:"get_sysinfo"`
	} `json:"system"`
	Emeter struct {
		GetRealtime struct {
			Current float64 `json:"current"`
			Voltage float64 `json:"voltage"`
			Power   float64 `json:"power"`
			Total   float64 `json:"total"`
			ErrCode int     `json:"err_code"`
		} `json:"get_realtime"`
		GetVgainIgain struct {
			Vgain   int `json:"vgain"`
			Igain   int `json:"igain"`
			ErrCode int `json:"err_code"`
		} `json:"get_vgain_igain"`
	} `json:"emeter"`
}

var tplink kasa

var (
	Pomvoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_voltage",
		Help: "Voltage recorded by TPlink HS110",
	})
)

var (
	Pomcurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_current",
		Help: "Current recorded by TPlink HS110",
	})
)

var (
	Pompower = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_power",
		Help: "Power recorded by TPlink HS110",
	})
)

var (
	Pomtotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_total",
		Help: "Total recorded by TPlink HS110",
	})
)

func main() {
	log.SetOutput(os.Stdout)
	log.Info("starting")
	log.SetFormatter(&log.JSONFormatter{})

	pReg := prometheus.NewRegistry()
	pReg.MustRegister(Pomvoltage)
	pReg.MustRegister(Pomcurrent)
	pReg.MustRegister(Pompower)
	pReg.MustRegister(Pomtotal)
	go func() {
		for {
			process()
			pomStats()
			time.Sleep(10 * time.Second)
		}
	}()
	handler := promhttp.HandlerFor(pReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	http.ListenAndServe(":8089", nil)
	log.Info("Serving on port 8089")
}

func pomStats() {
	voltage := tplink.Emeter.GetRealtime.Voltage
	current := tplink.Emeter.GetRealtime.Current
	power := tplink.Emeter.GetRealtime.Power
	total := tplink.Emeter.GetRealtime.Total
	Pomvoltage.Set(voltage)
	Pomcurrent.Set(current)
	Pompower.Set(power)
	log.WithFields(log.Fields{
		"Power":   power,
		"Current": current,
		"Voltage": voltage,
		"Total":   total})
}

func process() {
	plug := hs1xxplug.Hs1xxPlug{IPAddress: os.Getenv("TPLINK_ADDR")}
	results, err := plug.MeterInfo()
	if err != nil {
		log.Error("err:", err)
	}
	error := json.Unmarshal([]byte(results), &tplink)
	if error != nil {
		log.Error(err)
	}
}
