package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nerd2/neohubgo"
)

func main() {
	username := os.Args[1]
	password := os.Args[2]

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

	fmt.Printf("Found device '%s' and online\n", dev.DeviceName)

	data, err := nh.GetData(dev.DeviceId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, liveDev := range data.CacheValue.LiveInfo.Devices {
		log.Printf("%s: current: %s target: %s heating: %t\n", liveDev.ZoneName, liveDev.ActualTemp, liveDev.SetTemp, liveDev.HeatOn)
	}
}
