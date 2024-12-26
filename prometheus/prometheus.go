package main

import (
	"log"
	"net/http"
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
			[]string{"zoneName"}, nil,
		),
		targetTemp: prometheus.NewDesc("target_temp",
			"The target temperature of the zone",
			[]string{"zoneName"}, nil,
		),
		isHeating: prometheus.NewDesc("is_heating",
			"Whether the current zone is heating",
			[]string{"zoneName"}, nil,
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
		log.Fatalln(err.Error())
	}

	var dev *neohubgo.Device
	if len(devices) == 0 {
		log.Fatalln("No devices")
	} else if len(devices) > 1 {
		deviceNames := []string{}
		for _, device := range devices {
			if len(os.Args) >= 3 && os.Args[3] == device.DeviceName {
				dev = &device
				break
			}
			deviceNames = append(deviceNames, device.DeviceName)
		}
		if len(os.Args) < 4 {
			log.Fatalf("Supply device name, options: %s\n", strings.Join(deviceNames, ","))
		} else if dev == nil {
			log.Fatalf("Requested device name not found")
		}
	} else {
		dev = &devices[0]
	}

	if !dev.Online {
		log.Fatalln("Device offline")
	}

	data, err := nh.GetData(dev.DeviceId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, liveDev := range data.CacheValue.LiveInfo.Devices {
		actualTemp, _ := strconv.ParseFloat(liveDev.ActualTemp, 64)
		targetTemp, _ := strconv.ParseFloat(liveDev.SetTemp, 64)
		var isHeating float64
		if liveDev.HeatOn {
			isHeating = 1.0
		} else {
			isHeating = 0.0
		}

		ch <- prometheus.MustNewConstMetric(collector.currentTemp, prometheus.GaugeValue, actualTemp, liveDev.ZoneName)
		ch <- prometheus.MustNewConstMetric(collector.targetTemp, prometheus.GaugeValue, targetTemp, liveDev.ZoneName)
		ch <- prometheus.MustNewConstMetric(collector.isHeating, prometheus.GaugeValue, isHeating, liveDev.ZoneName)

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
