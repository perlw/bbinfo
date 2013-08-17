package main

import (
	"fmt"
	"github.com/conformal/gotk3/gtk"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
	"io/ioutil"
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

type NetworkStatus struct {
	Strength      int
	State         NetState
	Network       string
	LastTime      int
	ConnectedTime int
	CurrentUp     int
	CurrentDown   int
	TotalUp       int
	TotalDown     int
}

var modemOnline = false
var networkStatus = NetworkStatus{}

func bytesToHumanReadable(bytes int) string {
	switch {
	case bytes >= SizePB:
		return fmt.Sprintf("%.2fpb", float32(bytes)/float32(SizePB))
	case bytes >= SizeTB:
		return fmt.Sprintf("%.2ftb", float32(bytes)/float32(SizeTB))
	case bytes >= SizeGB:
		return fmt.Sprintf("%.2fgb", float32(bytes)/float32(SizeGB))
	case bytes >= SizeMB:
		return fmt.Sprintf("%.2fmb", float32(bytes)/float32(SizeMB))
	case bytes >= SizeKB:
		return fmt.Sprintf("%.2fkb", float32(bytes)/float32(SizeKB))
	default:
		return fmt.Sprintf("%db", bytes)
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
	networkStatus.Network = status[4]
	networkStatus.LastTime, _ = strconv.Atoi(status[7])
	networkStatus.ConnectedTime, _ = strconv.Atoi(status[9])
	networkStatus.CurrentDown, _ = strconv.Atoi(status[10])
	networkStatus.CurrentUp, _ = strconv.Atoi(status[11])
	networkStatus.TotalDown, _ = strconv.Atoi(status[13])
	networkStatus.TotalUp, _ = strconv.Atoi(status[14])

	fmt.Println(networkStatus)
	currentTransfer := networkStatus.CurrentDown + networkStatus.CurrentUp
	fmt.Printf("%s (%s/%s)\n", bytesToHumanReadable(currentTransfer), bytesToHumanReadable(networkStatus.CurrentDown), bytesToHumanReadable(networkStatus.CurrentUp))
}

func pollStatus(indicator *gotk3.AppIndicatorGotk3) {
	ticker := time.Tick(time.Second)

	for _ = range ticker {
		if response, err := http.Get("http://192.168.0.1/goform/status_update"); err != nil {
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
				strength := StrengthTable[networkStatus.Strength]
				fmt.Println(strength)
				fmt.Println("nm-signal-" + strength)
				indicator.SetIcon("nm-signal-"+strength, "Modem online")

				fmt.Println(string(data))
				parseStatusString(string(data))
			}
		}
	}
}

func main() {
	gtk.Init(nil)

	menu, _ := gtk.MenuNew()
	menuQuit, _ := gtk.MenuItemNewWithLabel("Quit")
	menuQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	menuQuit.Show()
	menu.Append(menuQuit)

	indicator := gotk3.NewAppIndicator("bb-4g-modem-indicator", "nm-signal-0", appindicator.CategoryCommunications)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	go pollStatus(indicator)

	gtk.Main()
}
