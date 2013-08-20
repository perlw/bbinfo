package main

import (
	"bbinfo/bytesconv"
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

var StrengthTable = []string{"0", "0", "25", "50", "75", "100"}
var RadioTable = []string{"-", "LTE", "EVD0", "CDMA1x", "WCDMA", "GSM", "HSUPA", "HSPA+", "DC-HSPA+", "EDGE", "GPRS"}

type NetworkStatus struct {
	Strength      int
	State         NetState
	RoamStatus    int
	Radio         int
	Network       string
	PINStatus     int
	ShowUnreadSMS int
	LastTime      int
	GetUnreadSMS  int
	ConnectedTime int
	CurrentUp     int
	CurrentDown   int
	TotalUp       int
	TotalDown     int
	ServiceStatus int
	SpeedDown     int
	SpeedUp       int
}

var modemOnline = false
var networkStatus = NetworkStatus{}

func timestampToString(timestamp int) string {
	hour := timestamp / 3000
	min := timestamp / 60
	sec := timestamp - (hour * 360) - (min * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
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
	networkStatus.RoamStatus, _ = strconv.Atoi(status[2])
	networkStatus.Radio, _ = strconv.Atoi(status[3])
	networkStatus.Network = status[4]
	networkStatus.PINStatus, _ = strconv.Atoi(status[5])
	networkStatus.ShowUnreadSMS, _ = strconv.Atoi(status[6])
	networkStatus.LastTime, _ = strconv.Atoi(status[7])
	networkStatus.GetUnreadSMS, _ = strconv.Atoi(status[8])
	networkStatus.ConnectedTime, _ = strconv.Atoi(status[9])
	networkStatus.CurrentDown, _ = strconv.Atoi(status[10])
	networkStatus.CurrentUp, _ = strconv.Atoi(status[11])
	networkStatus.TotalDown, _ = strconv.Atoi(status[13])
	networkStatus.TotalUp, _ = strconv.Atoi(status[14])
	networkStatus.ServiceStatus, _ = strconv.Atoi(status[15])
	networkStatus.SpeedDown, _ = strconv.Atoi(status[16])
	networkStatus.SpeedUp, _ = strconv.Atoi(status[17])

	if networkStatus.Radio > -1 {
		fmt.Printf("Radio: %s\n", RadioTable[networkStatus.Radio])
	}
	fmt.Printf("Network: %s\n", networkStatus.Network)
	fmt.Printf("PINStatus: %d\n", networkStatus.PINStatus)
	//fmt.Printf("LastTime: %s\n", time.Unix(int64(networkStatus.LastTime), 0).String())
	fmt.Printf("ConnectedTime: %s\n", timestampToString(networkStatus.ConnectedTime))
	fmt.Printf("CurrentDown: %s\n", bytesconv.ToHumanReadable(networkStatus.CurrentDown))
	fmt.Printf("CurrentUp: %s\n", bytesconv.ToHumanReadable(networkStatus.CurrentUp))
	fmt.Printf("TotalDown: %s\n", bytesconv.ToHumanReadable(networkStatus.TotalDown))
	fmt.Printf("TotalUp: %s\n", bytesconv.ToHumanReadable(networkStatus.TotalUp))
	fmt.Printf("ServiceStatus: %d\n", networkStatus.ServiceStatus)
	fmt.Printf("SpeedDown: %s\n", bytesconv.ToHumanReadable(networkStatus.SpeedDown))
	fmt.Printf("SpeedUp: %s\n", bytesconv.ToHumanReadable(networkStatus.SpeedUp))
}

func pollStatus(indicator *gotk3.AppIndicatorGotk3, menuCurrent, menuTotal *gtk.MenuItem) {
	ticker := time.Tick(time.Second)

	for _ = range ticker {
		req := fmt.Sprintf("http://192.168.0.1/goform/status_update?status_flag=1&rd=%f", rand.Float32())
		if response, err := http.Get(req); err != nil {
			if modemOnline {
				modemOnline = false
				indicator.SetIcon("network-error", "Modem offline")
				indicator.SetLabel("", "")
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
				if networkStatus.Radio > -1 {
					indicator.SetLabel(timestampToString(networkStatus.ConnectedTime), "")
				} else {
					indicator.SetLabel("", "")
				}

				<-time.After(time.Millisecond * 10)
				cUp, cDown, cQual := bytesconv.QualifyTransfer(networkStatus.CurrentUp, networkStatus.CurrentDown)
				currentStr := fmt.Sprintf("Current Transfer - %.2f/%.2f%s", cUp, cDown, cQual)
				menuCurrent.SetLabel(currentStr)

				<-time.After(time.Millisecond * 10)
				tUp, tDown, tQual := bytesconv.QualifyTransfer(networkStatus.TotalUp, networkStatus.TotalDown)
				totalStr := fmt.Sprintf("Total Transfer - %.2f/%.2f%s", tUp, tDown, tQual)
				menuTotal.SetLabel(totalStr)
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

	indicator := gotk3.NewAppIndicator("bb-4g-modem-indicator", "network-error", appindicator.CategoryCommunications)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	go pollStatus(indicator, menuCurrentTransfer, menuTotalTransfer)

	gtk.Main()
}
