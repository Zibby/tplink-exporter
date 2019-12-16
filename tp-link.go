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

// RequestedPlug holds info about hs110 plugs
type RequestedPlug struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Legacy   bool   `json:"legacy"`
	Handler  http.Handler
	Registry prometheus.Registry
	Stats    struct {
		Voltage float64
		Current float64
		Power   float64
	}
}

type kasaNew struct {
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

type kasaOld struct {
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

var plugsIn []RequestedPlug

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

//var plugIP = os.Getenv("TPLINK_ADDR")
//var plug = hs1xxplug.Hs1xxPlug{IPAddress: plugIP}

func register(preg *prometheus.Registry) {
	log.Info("Registering Stats")
	preg.MustRegister(Pomvoltage)
	preg.MustRegister(Pomcurrent)
	preg.MustRegister(Pompower)
	preg.MustRegister(Pomtotal)
	log.Info("Stats Registered")
}

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
}

func decodeEnvJSON() {
	log.Info("Decoding Env Json")
	err := json.Unmarshal([]byte(os.Getenv("PLUGS")), &plugsIn)
	if err != nil {
		log.Info("Error with Env Var PLUGS:", err)
	}
	log.Info("Plugs Requested: ", &plugsIn)
}

func init() {
	initLog()
	decodeEnvJSON()
	go func() {
		http.ListenAndServe(":8089", nil)
		log.Info("Serving on port 8089")
	}()
}

func runloop(rPlug RequestedPlug) {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			log.Error("Timed out loop")
		case <-tick:
			connectedToPlug := connectToPlug(rPlug)
			if connectedToPlug == true {
				pomStats(rPlug)
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func main() {
	log.Info(plugsIn)
	for _, plug := range plugsIn {
		plug.Registry = *prometheus.NewRegistry()
		plug.Handler = promhttp.HandlerFor(&plug.Registry, promhttp.HandlerOpts{})
	}
	for _, plug := range plugsIn {
		log.Info("Range: ", plug.Name)
		go func(plug RequestedPlug) {
			log.Info("Handling: ", plug.Name)
			http.Handle("/"+plug.Name, plug.Handler)
			runloop(plug)
		}(plug)
		time.Sleep(10 * time.Second)
	}
}

func pomStats(rplug RequestedPlug) {
	if rplug.Legacy == true {
		Pomvoltage.Set(rplug.Stats.Voltage)
		Pomcurrent.Set(rplug.Stats.Current)
		Pompower.Set(rplug.Stats.Power)
		log.WithFields(log.Fields{
			"Power":   rplug.Stats.Power,
			"Current": rplug.Stats.Current,
			"Voltage": rplug.Stats.Voltage,
		}).Info("Publishing Stats for:", rplug.Name)
	} else {
		log.WithFields(log.Fields{
			"Power":   rplug.Stats.Power,
			"Current": rplug.Stats.Current,
			"Voltage": rplug.Stats.Voltage,
		}).Info("Publishing Stats")
	}
}

func connectToPlug(rplug RequestedPlug) bool {
	log.WithFields(log.Fields{
		"Plug_IP": rplug.Address,
	}).Info("connecting to plug")
	h1plug := hs1xxplug.Hs1xxPlug{IPAddress: rplug.Address}
	log.Info("h1plug set")
	readings, err := h1plug.MeterInfo()
	log.Info("Got Readings")
	if err != nil {
		log.Error(err)
		return false
	}
	log.Info("Unmarshaling meter reading")
	if rplug.Legacy == false {
		log.Info("Using Later FW json")
		var results kasaNew
		err = json.Unmarshal([]byte(readings), &results)
		rplug.Stats.Voltage = results.Emeter.GetRealtime.Voltage / 1000
		rplug.Stats.Current = results.Emeter.GetRealtime.Current / 1000
		rplug.Stats.Power = results.Emeter.GetRealtime.Power / 1000
		log.Info(rplug.Stats.Voltage)
	} else {
		log.Info("Using legacy FW json")
		var results kasaOld
		err = json.Unmarshal([]byte(readings), &results)
		rplug.Stats.Voltage = results.Emeter.GetRealtime.Voltage
		rplug.Stats.Current = results.Emeter.GetRealtime.Current
		rplug.Stats.Power = results.Emeter.GetRealtime.Power
		log.Info(rplug.Stats.Voltage)

		return false
	}

	if err != nil {
		log.Info(err)
		return false
	}
	return false
}
