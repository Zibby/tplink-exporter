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

type kasa_new struct {
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
			Current float64 `json:"current_ma"`
			Voltage float64 `json:"voltage_mv"`
			Power   float64 `json:"power_mw"`
			Total   float64 `json:"total_wh"`
			ErrCode int     `json:"err_code"`
		} `json:"get_realtime"`
		GetVgainIgain struct {
			Vgain   int `json:"vgain"`
			Igain   int `json:"igain"`
			ErrCode int `json:"err_code"`
		} `json:"get_vgain_igain"`
	} `json:"emeter"`
}

type kasa_old struct {
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



var tplinkOld kasa_old
var tplinkNew kasa_new

var (
	// Pomvoltage is the current voltage will see
	Pomvoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_voltage",
		Help: "Voltage recorded by TPlink HS110",
	})
)

var (
	// Pomcurrent is the current prometheus will see
	Pomcurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_current",
		Help: "Current recorded by TPlink HS110",
	})
)

var (
	// Pompower is the power prometheus will see
	Pompower = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_power",
		Help: "Power recorded by TPlink HS110",
	})
)

var (
	// Pomtotal is the total prometheus will see
	Pomtotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kasa_total",
		Help: "Total recorded by TPlink HS110",
	})
)
var plugIP = os.Getenv("TPLINK_ADDR")
var plug = hs1xxplug.Hs1xxPlug{IPAddress: plugIP}
var pReg = prometheus.NewRegistry()

func register() {
	log.Info("Registering Stats")
	pReg.MustRegister(Pomvoltage)
	pReg.MustRegister(Pomcurrent)
	pReg.MustRegister(Pompower)
	pReg.MustRegister(Pomtotal)
	log.Info("Stats Registered")
}

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
}

func init() {
	initLog()
	register()
}

func serve() {
	log.Info("starting http hander")
	handler := promhttp.HandlerFor(pReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	http.ListenAndServe(":8089", nil)
	log.Info("Serving on port 8089")
}

func runloop() {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			log.Error("Timed out loop")
		case <-tick:
			connectedToPlug := connectToPlug()
			if connectedToPlug == true {
				pomStats()
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func main() {
	go func() {
		serve()
	}()
	runloop()
}

func pomStats() {
	if os.Getenv("LATER_FW") == "true"{
	  voltage := tplinkNew.Emeter.GetRealtime.Voltage/1000
	  current := tplinkNew.Emeter.GetRealtime.Current/1000
	  power := tplinkNew.Emeter.GetRealtime.Power/1000
	  total := tplinkNew.Emeter.GetRealtime.Total/1000
	  Pomvoltage.Set(voltage)
	  Pomcurrent.Set(current)
	  Pompower.Set(power)
	  Pomtotal.Set(total)
	  log.WithFields(log.Fields{
		  "Power":   power,
		  "Current": current,
		  "Voltage": voltage,
		  "Total":   total,
	  }).Info("Publishing Stats")
	} else {
	  voltage := tplinkOld.Emeter.GetRealtime.Voltage
	  current := tplinkOld.Emeter.GetRealtime.Current
	  power := tplinkOld.Emeter.GetRealtime.Power
	  total := tplinkOld.Emeter.GetRealtime.Total
	  Pomvoltage.Set(voltage)
	  Pomcurrent.Set(current)
	  Pompower.Set(power)
	  Pomtotal.Set(total)
	  log.WithFields(log.Fields{
		  "Power":   power,
		  "Current": current,
		  "Voltage": voltage,
		  "Total":   total,
	  }).Info("Publishing Stats")
	} 
}

func connectToPlug() bool {
	log.WithFields(log.Fields{
		"Plug_IP": plugIP,
	}).Info("connecting to plug")
	results, err := plug.MeterInfo()
	if err != nil {
		log.Error(err)
		return false
	}
	log.Info("Unmarshaling meter reading")
	if os.Getenv("LATER_FW") == "true" {
		err = json.Unmarshal([]byte(results), &tplinkNew)
	} else {
		err = json.Unmarshal([]byte(results), &tplinkOld)
	}
	
	if err != nil {
		log.Info(err)
		return false
	}
	return true
}
