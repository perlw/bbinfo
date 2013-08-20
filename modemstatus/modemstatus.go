package modemstatus

import (
	"bbinfo/bytesconv"
	"fmt"
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

type Status struct {
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

func (s Status) ToString() string {
	str := ""
	if s.Radio > -1 {
		str += fmt.Sprintf("Radio: %s\n", RadioTable[s.Radio])
	}
	str += fmt.Sprintf("Network: %s\n", s.Network)
	str += fmt.Sprintf("PINStatus: %d\n", s.PINStatus)
	//str += fmt.Sprintf("LastTime: %s\n", time.Unix(int64(s.LastTime), 0).String())
	str += fmt.Sprintf("ConnectedTime: %s\n", TimestampToString(s.ConnectedTime))
	str += fmt.Sprintf("CurrentDown: %s\n", bytesconv.ToHumanReadable(s.CurrentDown))
	str += fmt.Sprintf("CurrentUp: %s\n", bytesconv.ToHumanReadable(s.CurrentUp))
	str += fmt.Sprintf("TotalDown: %s\n", bytesconv.ToHumanReadable(s.TotalDown))
	str += fmt.Sprintf("TotalUp: %s\n", bytesconv.ToHumanReadable(s.TotalUp))
	str += fmt.Sprintf("ServiceStatus: %d\n", s.ServiceStatus)
	str += fmt.Sprintf("SpeedDown: %s\n", bytesconv.ToHumanReadable(s.SpeedDown))
	str += fmt.Sprintf("SpeedUp: %s\n", bytesconv.ToHumanReadable(s.SpeedUp))

	return str
}

//5;2;0;9;Telenor SE;1;;;;797;4044181;1603028;59801;1274938731;83429417;2;3608;4848;
func parseStatusString(data string) *Status {
	status := strings.Split(data, ";")

	modemStatus := Status{}
	modemStatus.Strength, _ = strconv.Atoi(status[0])
	netState, _ := strconv.Atoi(status[1])
	if netState == 2 {
		modemStatus.State = StateConnected
	} else {
		modemStatus.State = StateDisconnected
	}
	modemStatus.RoamStatus, _ = strconv.Atoi(status[2])
	modemStatus.Radio, _ = strconv.Atoi(status[3])
	modemStatus.Network = status[4]
	modemStatus.PINStatus, _ = strconv.Atoi(status[5])
	modemStatus.ShowUnreadSMS, _ = strconv.Atoi(status[6])
	modemStatus.LastTime, _ = strconv.Atoi(status[7])
	modemStatus.GetUnreadSMS, _ = strconv.Atoi(status[8])
	modemStatus.ConnectedTime, _ = strconv.Atoi(status[9])
	modemStatus.CurrentDown, _ = strconv.Atoi(status[10])
	modemStatus.CurrentUp, _ = strconv.Atoi(status[11])
	modemStatus.TotalDown, _ = strconv.Atoi(status[13])
	modemStatus.TotalUp, _ = strconv.Atoi(status[14])
	modemStatus.ServiceStatus, _ = strconv.Atoi(status[15])
	modemStatus.SpeedDown, _ = strconv.Atoi(status[16])
	modemStatus.SpeedUp, _ = strconv.Atoi(status[17])

	return &modemStatus
}

func TimestampToString(timestamp int) string {
	hour := timestamp / 3000
	min := timestamp / 60
	sec := timestamp - (hour * 360) - (min * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func DoPoll(tickFunc func(*Status), errFunc func(error)) {
	ticker := time.Tick(time.Second)

	for _ = range ticker {
		req := fmt.Sprintf("http://192.168.0.1/goform/status_update?status_flag=1&rd=%f", rand.Float32())
		if response, err := http.Get(req); err != nil {
			errFunc(err)
			fmt.Println("Could not connect, modem not available?")
		} else {
			defer response.Body.Close()
			if data, err := ioutil.ReadAll(response.Body); err != nil {
				fmt.Println("Error occurred while reading data.")
			} else {
				status := parseStatusString(string(data))
				tickFunc(status)
			}
		}
	}
}
