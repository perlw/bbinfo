package main

import (
	"fmt"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
	"github.com/doxxan/bbinfo/bytesconv"
	"github.com/doxxan/bbinfo/modemstatus"
	"github.com/doxxan/gotk3/gtk"
	"time"
)

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

	go func() {
		modemstatus.DoPoll(func(status *modemstatus.Status) {
			// Debug
			//fmt.Println(status.ToString())

			strength := modemstatus.StrengthTable[status.Strength]
			indicator.SetIcon("nm-signal-"+strength, "Modem online")
			if status.Radio > -1 {
				indicator.SetLabel(modemstatus.TimestampToString(status.ConnectedTime), "")
			} else {
				indicator.SetLabel("", "")
			}

			<-time.After(time.Millisecond * 10)
			cUp, cDown, cQual := bytesconv.QualifyTransfer(status.CurrentUp, status.CurrentDown)
			currentStr := fmt.Sprintf("Current Transfer - %.2f/%.2f%s", cUp, cDown, cQual)
			menuCurrentTransfer.SetLabel(currentStr)

			<-time.After(time.Millisecond * 10)
			tUp, tDown, tQual := bytesconv.QualifyTransfer(status.TotalUp, status.TotalDown)
			totalStr := fmt.Sprintf("Total Transfer - %.2f/%.2f%s", tUp, tDown, tQual)
			menuTotalTransfer.SetLabel(totalStr)
		}, func(err error) {
			indicator.SetIcon("network-error", "Modem offline")
			indicator.SetLabel("", "")
		})
	}()

	gtk.Main()
}
