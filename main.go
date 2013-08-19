package main

import (
	"fmt"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
	"github.com/doxxan/gotk3/gtk"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NetState int

const (
	StateDisconnected NetState = iota
	StateConnected
)

const (
	SizeKB = 1024
	SizeMB = SizeKB * 1024
	SizeGB = SizeMB * 1024
	SizeTB = SizeGB * 1024
	SizePB = SizeTB * 1024
)

var StrengthTable = []string{"0", "0", "25", "50", "75", "100"}
var RadioTable = []string{"-", "LTE", "EVD0", "CDMA1x", "WCDMA", "GSM", "HSUPA", "HSPA+", "DC-HSPA+", "EDGE", "GPRS"}

type NetworkStatus struct {
	Strength      int
	State         NetState
	Radio         int
	Network       string
	PINStatus     int
	LastTime      int
	ConnectedTime int
	CurrentUp     int
	CurrentDown   int
	TotalUp       int
	TotalDown     int
	ServiceStatus int
}

var modemOnline = false
var networkStatus = NetworkStatus{}

func bytesToHumanReadable(bytes int) string {
	switch {
	case bytes >= SizePB:
		return fmt.Sprintf("%.2fPB", float32(bytes)/float32(SizePB))
	case bytes >= SizeTB:
		return fmt.Sprintf("%.2fTB", float32(bytes)/float32(SizeTB))
	case bytes >= SizeGB:
		return fmt.Sprintf("%.2fGB", float32(bytes)/float32(SizeGB))
	case bytes >= SizeMB:
		return fmt.Sprintf("%.2fMB", float32(bytes)/float32(SizeMB))
	case bytes >= SizeKB:
		return fmt.Sprintf("%.2fKB", float32(bytes)/float32(SizeKB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

//5;2;0;9;Telenor SE;1;;;;797;4044181;1603028;59801;1274938731;83429417;2;3608;4848;
func parseStatusString(data string) {
	status := strings.Split(data, ";")

	networkStatus.Strength, _ = strconv.Atoi(status[0])
	netState, _ := strconv.Atoi(status[1])
	if netState == 2 {
		networkStatus.State = StateConnected
	} else {
		networkStatus.State = StateDisconnected
	}
	networkStatus.Radio, _ = strconv.Atoi(status[3])
	networkStatus.Network = status[4]
	networkStatus.PINStatus, _ = strconv.Atoi(status[5])
	networkStatus.LastTime, _ = strconv.Atoi(status[7])
	networkStatus.ConnectedTime, _ = strconv.Atoi(status[9])
	networkStatus.CurrentDown, _ = strconv.Atoi(status[10])
	networkStatus.CurrentUp, _ = strconv.Atoi(status[11])
	networkStatus.TotalDown, _ = strconv.Atoi(status[13])
	networkStatus.TotalUp, _ = strconv.Atoi(status[14])
	networkStatus.ServiceStatus, _ = strconv.Atoi(status[15])

	if networkStatus.Radio > -1 {
		fmt.Printf("Radio: %s\n", RadioTable[networkStatus.Radio])
	}
	fmt.Printf("Network: %s\n", networkStatus.Network)
	fmt.Printf("PINStatus: %d\n", networkStatus.PINStatus)
	//fmt.Printf("LastTime: %s\n", time.Unix(int64(networkStatus.LastTime), 0).String())
	hour := networkStatus.ConnectedTime / 3000
	min := networkStatus.ConnectedTime / 60
	sec := networkStatus.ConnectedTime - (hour * 360) - (min * 60)
	fmt.Printf("ConnectedTime: %d:%d:%d\n", hour, min, sec)
	fmt.Printf("CurrentDown: %s\n", bytesToHumanReadable(networkStatus.CurrentDown))
	fmt.Printf("CurrentUp: %s\n", bytesToHumanReadable(networkStatus.CurrentUp))
	fmt.Printf("TotalDown: %s\n", bytesToHumanReadable(networkStatus.TotalDown))
	fmt.Printf("TotalUp: %s\n", bytesToHumanReadable(networkStatus.TotalUp))
	fmt.Printf("ServiceStatus: %d\n", networkStatus.ServiceStatus)
}

func pollStatus(indicator *gotk3.AppIndicatorGotk3, menuCurrent, menuTotal *gtk.MenuItem) {
	ticker := time.Tick(time.Second)

	either := false
	for _ = range ticker {
		req := fmt.Sprintf("http://192.168.0.1/goform/status_update?rd=%f", rand.Float32())
		if response, err := http.Get(req); err != nil {
			if modemOnline {
				modemOnline = false
				indicator.SetIcon("network-error", "Modem offline")
			}
			fmt.Println("Could not connect, modem not available?")
		} else {
			defer response.Body.Close()
			if data, err := ioutil.ReadAll(response.Body); err != nil {
				fmt.Println("Error occurred while reading data.")
			} else {
				if !modemOnline {
					modemOnline = true
				}
				fmt.Println(string(data))
				parseStatusString(string(data))

				strength := StrengthTable[networkStatus.Strength]
				indicator.SetIcon("nm-signal-"+strength, "Modem online")

				if either {
					currentStr := fmt.Sprintf("Current Transfer - U:%s | D:%s", bytesToHumanReadable(networkStatus.CurrentUp), bytesToHumanReadable(networkStatus.CurrentDown))
					menuCurrent.SetLabel(currentStr)
				} else {
					totalStr := fmt.Sprintf("Total Transfer - U:%s | D:%s", bytesToHumanReadable(networkStatus.TotalUp), bytesToHumanReadable(networkStatus.TotalDown))
					menuTotal.SetLabel(totalStr)
				}
				either = !either
			}
		}
	}
}

func main() {
	gtk.Init(nil)

	menu, _ := gtk.MenuNew()

	menuCurrentTransfer, _ := gtk.MenuItemNewWithLabel("Current Transfer - ##")
	menuCurrentTransfer.Show()
	menu.Append(menuCurrentTransfer)

	menuTotalTransfer, _ := gtk.MenuItemNewWithLabel("Total Transfer - ##")
	menuTotalTransfer.Show()
	menu.Append(menuTotalTransfer)

	menuQuit, _ := gtk.MenuItemNewWithLabel("Quit")
	menuQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	menuQuit.Show()
	menu.Append(menuQuit)

	indicator := gotk3.NewAppIndicator("bb-4g-modem-indicator", "nm-signal-0", appindicator.CategoryCommunications)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	go pollStatus(indicator, menuCurrentTransfer, menuTotalTransfer)

	gtk.Main()
}
