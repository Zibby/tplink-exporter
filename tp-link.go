package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sausheong/hs1xxplug"
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

func init() {
	prometheus.MustRegister(Pomvoltage)
	prometheus.MustRegister(Pomcurrent)
	prometheus.MustRegister(Pompower)

}

func main() {
	go func() {
		for {
			process()
			pomStats()
			time.Sleep(10 * time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8089", nil)
}

func pomStats() {
	Pomvoltage.Set(tplink.Emeter.GetRealtime.Voltage)
	Pomcurrent.Set(tplink.Emeter.GetRealtime.Current)
	Pompower.Set(tplink.Emeter.GetRealtime.Power)
}

func process() {
	plug := hs1xxplug.Hs1xxPlug{IPAddress: os.Getenv("TPLINK_ADDR")}
	results, err := plug.MeterInfo()
	if err != nil {
		fmt.Println("err:", err)
	}
	error := json.Unmarshal([]byte(results), &tplink)
	if error != nil {
		fmt.Println(err)
	}
}
