package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nerd2/neohubgo"
)

type neoCollector struct {
	currentTemp *prometheus.Desc
	targetTemp  *prometheus.Desc
	isHeating   *prometheus.Desc
}

func newNeoCollector() *neoCollector {
	return &neoCollector{
		currentTemp: prometheus.NewDesc("current_temp",
			"The current temperature of the zone",
			[]string{"zone", "device"}, nil,
		),
		targetTemp: prometheus.NewDesc("target_temp",
			"The target temperature of the zone",
			[]string{"zone", "device"}, nil,
		),
		isHeating: prometheus.NewDesc("is_heating",
			"Whether the current zone is heating",
			[]string{"zone", "device"}, nil,
		),
	}
}

func (collector *neoCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.currentTemp
	ch <- collector.targetTemp
	ch <- collector.isHeating
}

// Collect implements required collect function for all promehteus collectors
func (collector *neoCollector) Collect(ch chan<- prometheus.Metric) {

	username := os.Getenv("NEOPROM_USERNAME")
	password := os.Getenv("NEOPROM_PASSWORD")

	nh := neohubgo.NewNeoHub(&neohubgo.Options{Username: username, Password: password})
	devices, err := nh.Login()
	if err != nil {
		log.Println("Error : couldn't connect -> ", err)
		return
	}

	for _, device := range devices {
		if device.Online {
			deviceName, _ := url.QueryUnescape(strings.Trim(device.DeviceName, " "))

			data, err := nh.GetData(device.DeviceId)
			if err != nil {
				// If we've got this far, this may be recoverable and we may be able to query other devices
				log.Println("Warn : Unable to query device -> ", err)
			}

			for _, liveDev := range data.CacheValue.LiveInfo.Devices {
				actualTemp, _ := strconv.ParseFloat(liveDev.ActualTemp, 64)
				targetTemp, _ := strconv.ParseFloat(liveDev.SetTemp, 64)
				zoneName, _ := url.QueryUnescape(strings.Trim(liveDev.ZoneName, " "))

				var isHeating float64
				if liveDev.HeatOn {
					isHeating = 1.0
				} else {
					isHeating = 0.0
				}

				ch <- prometheus.MustNewConstMetric(collector.currentTemp, prometheus.GaugeValue, actualTemp, zoneName, deviceName)
				ch <- prometheus.MustNewConstMetric(collector.targetTemp, prometheus.GaugeValue, targetTemp, zoneName, deviceName)
				ch <- prometheus.MustNewConstMetric(collector.isHeating, prometheus.GaugeValue, isHeating, zoneName, deviceName)

			}
		}
	}
}

func main() {
	neo := newNeoCollector()
	prometheus.MustRegister(neo)
	port := "9101"
	log.Printf("Starting metrics server on port %s", port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
