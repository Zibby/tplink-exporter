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

var plugsIn []RequestedPlug

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

func worker(job <-chan RequestedPlug) {
	for {
		for j := range job {
			reg := prometheus.NewRegistry()
			handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
			var (
				promTotal   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kasa_total", Help: "Total recorded by TPlink HS110"})
				promPower   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kasa_power", Help: "Power recorded by TPlink HS110"})
				promCurrent = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kasa_current", Help: "Current recorded by TPlink HS110"})
				promVoltage = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kasa_voltage", Help: "Voltage recorded by TPlink HS110"})
			)
			reg.MustRegister(promVoltage)
			reg.MustRegister(promCurrent)
			reg.MustRegister(promPower)
			reg.MustRegister(promTotal)
			http.Handle("/"+j.Name, handler)
			log.Info("Handling: ", j.Name)
			h1plug := hs1xxplug.Hs1xxPlug{IPAddress: j.Address}
			for {
				timeout := time.After(5 * time.Second)
				tick := time.Tick(500 * time.Millisecond)
				select {
				case <-timeout:
					log.Error("Timed out loop")
				case <-tick:
					readings, err := h1plug.MeterInfo()
					if err != nil {
						log.Error(err)
					}
					log.Info("Unmarshaling meter reading")
					if j.Legacy == false {
						log.Info("Using Later FW json")
						var results KasaNew
						err = json.Unmarshal([]byte(readings), &results)
						promVoltage.Set(results.Emeter.GetRealtime.Voltage / 1000)
						promCurrent.Set(results.Emeter.GetRealtime.Current / 1000)
						promPower.Set(results.Emeter.GetRealtime.Power / 1000)
					} else {
						log.Info("Using legacy FW json")
						var results kasaOld
						err = json.Unmarshal([]byte(readings), &results)
						promVoltage.Set(results.Emeter.GetRealtime.Voltage)
						promCurrent.Set(results.Emeter.GetRealtime.Current)
						promPower.Set(results.Emeter.GetRealtime.Power)
					}
				}
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func main() {
	jobs := make(chan RequestedPlug)
	log.Info(plugsIn)
	for _, plug := range plugsIn {
		go worker(jobs)
		jobs <- plug
	}
	for {
		time.Sleep(5 * time.Second)
	}
}
