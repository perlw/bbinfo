package main

import (
	"fmt"
	"github.com/conformal/gotk3/gtk"
	"github.com/perlw/appindicator"
	"github.com/perlw/appindicator/gtk-extensions/gotk3"
	"io/ioutil"
	"net/http"
	"time"
)

func checkStats() {
	ticker := time.Tick(time.Second)

	for _ = range ticker {
		if response, err := http.Get("http://192.168.0.1/goform/status_update"); err != nil {
			fmt.Println("Could not connect, modem not available?")
		} else {
			defer response.Body.Close()
			if data, err := ioutil.ReadAll(response.Body); err != nil {
				fmt.Println("Error occurred while reading data.")
			} else {
				fmt.Println(string(data))
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

	indicator := gotk3.NewAppIndicator("test-indicator", "indicator-messages", appindicator.CategoryApplicationStatus)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	go checkStats()

	gtk.Main()
}
