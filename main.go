package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sausheong/hs1xxplug"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
}

func init() {
	initLog()
}

func plugStats(w http.ResponseWriter, r *http.Request) {
	var (
		promTotal = prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "kasa_total", Help: "Total recorded by TPlink HS110"})
		promPower = prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "kasa_power", Help: "Power recorded by TPlink HS110"})
		promCurrent = prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "kasa_current", Help: "Current recorded by TPlink HS110"})
		promVoltage = prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "kasa_voltage", Help: "Voltage recorded by TPlink HS110"})
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(promVoltage)
	reg.MustRegister(promCurrent)
	reg.MustRegister(promPower)
	reg.MustRegister(promTotal)
	r.ParseForm()

	plugAddress := r.FormValue("address")
	useLegacy := r.FormValue("legacy")
	h1plug := hs1xxplug.Hs1xxPlug{IPAddress: plugAddress}
	readings, err := h1plug.MeterInfo()
	if err != nil {
		log.Error(err)
	}
	if useLegacy == "false" {
		log.Info("Using Later FW json")
		var results KasaNew
		err = json.Unmarshal([]byte(readings), &results)
		promVoltage.Set(results.Emeter.GetRealtime.Voltage / 1000)
		promCurrent.Set(results.Emeter.GetRealtime.Current / 1000)
		promPower.Set(results.Emeter.GetRealtime.Power / 1000)
		promPower.Set(results.Emeter.GetRealtime.Total)
	} else {
		log.Info("Using legacy FW json")
		var results kasaOld
		err = json.Unmarshal([]byte(readings), &results)
		promVoltage.Set(results.Emeter.GetRealtime.Voltage)
		promCurrent.Set(results.Emeter.GetRealtime.Current)
		promPower.Set(results.Emeter.GetRealtime.Power)
		promPower.Set(results.Emeter.GetRealtime.Total)
	}
	log.Info(promVoltage)
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func healthHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health Ok!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHander)
	r.HandleFunc("/metrics", plugStats)

	log.Fatal(http.ListenAndServe(":8089", r))
}
