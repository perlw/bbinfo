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
	StateConnected    NetState = iota
)

type NetworkStatus struct {
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

//5;2;0;9;Telenor SE;1;;;;797;4044181;1603028;59801;1274938731;83429417;2;3608;4848;
func parseStatusString(data string) {
	status := strings.Split(data, ";")

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
	//networkStatus.CurrentTransfer = networkStatus.CurrentDown + networkStatus.CurrentUp
	networkStatus.TotalDown, _ = strconv.Atoi(status[13])
	networkStatus.TotalUp, _ = strconv.Atoi(status[14])
	//networkStatus.TotalTransfer = networkStatus.TotalDown + networkStatus.TotalUp

	fmt.Printf(networkStatus)
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
					indicator.SetIcon("network-idle", "Modem online")
				}

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

	indicator := gotk3.NewAppIndicator("bb-4g-modem-indicator", "network-offline", appindicator.CategoryCommunications)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	go pollStatus(indicator)

	gtk.Main()
}

/*
var pin_status = 'SIM_READY';
var pin_retry = '3';

function initTranslation()
{
  var   e = document.getElementById("Current_flux");
		e.innerHTML = Current_flux[lang_index];
		e = document.getElementById("Total_flux");
		e.innerHTML = Total_flux[lang_index];
		e = document.getElementById("Connected_time");
		e.innerHTML = Connected_time[lang_index];
}

var Auto = '0';
function show_net_state(NetStat)
{
	if(NetStat=='2')
	{
		NetStat="Disconnect";
		//manual_mode_value = "Disconnect";
		document.getElementById("ppp_state_img").src="../img2/connect.png";
		//document.getElementById("manual_mode_show").value = Disconnect[lang_index];
		document.getElementById("manual_mode_show_sp").innerHTML = "<input id=\"manual_mode_show\" type=\"submit\" value=\"" +Disconnect[lang_index]+"\"   onClick=\"return connect_net();\" name=\"manual_mode_show\" class=\"homebtn\" onMouseOver=\"this.className=\'homemenu_over\'\" onMouseOut=\"this.className=\'homemenu_out\'\" onMouseDown=\"this.className=\'homemenu_down\'\">"
	}
	else if(NetStat=='0')
	{
		NetStat="Connect";
		//manual_mode_value = "Connect";
		document.getElementById("ppp_state_img").src="../img2/disconnect.png";
		//document.getElementById("manual_mode_show").value = Connect[lang_index];
		document.getElementById("manual_mode_show_sp").innerHTML = "<input id=\"manual_mode_show\" type=\"submit\" value=\"" +Connect[lang_index]+"\"   onClick=\"return connect_net();\" name=\"manual_mode_show\" class=\"homebtn\" onMouseOver=\"this.className=\'homemenu_over\'\" onMouseOut=\"this.className=\'homemenu_out\'\" onMouseDown=\"this.className=\'homemenu_down\'\">"
	}
	else
	{
		NetStat="Connect";
		document.getElementById("ppp_state_img").src="../img2/connecting.gif";
		//document.getElementById("manual_mode_show").value = Connect[lang_index];
		document.getElementById("manual_mode_show_sp").innerHTML = "<input id=\"manual_mode_show\" type=\"submit\" value=\"" +Connect[lang_index]+"\"   onClick=\"return connect_net();\" name=\"manual_mode_show\" class=\"homebtn\" onMouseOver=\"this.className=\'homemenu_over\'\" onMouseOut=\"this.className=\'homemenu_out\'\" onMouseDown=\"this.className=\'homemenu_down\'\">"
		window.location = "../WanSettings/connect_wait2_home_wan.asp";
	}
	document.getElementById("manual_mode").value = NetStat;

	if(NetStat != 'Disconnect')
	{
		if(('1' == '1')||('' == '1')||('0' == Auto))
		{
			document.getElementById('manual_mode_show').disabled = false;
		}
	}
	else
	{
		document.getElementById('manual_mode_show').disabled = false;
	}
}
function show_connect_time(ConnectNetTime,NetStat)
{
    if(NetStat == '0')
	{
	    document.getElementById("last_time").innerHTML = "00:00:00";
		return;
	}
	second=ConnectNetTime%60;
	if(second < 10)
	{
		second = '0' + second;
	}
	minite=parseInt((ConnectNetTime/60)%60);
	if(minite < 10)
	{
		minite = '0' + minite;
	}
	hour=parseInt(ConnectNetTime/3600);
	if(hour < 10)
	{
		hour = '0' + hour;
	}
	document.getElementById("last_time").innerHTML = hour+":"+minite+":"+second;
}
function show_decimal(decimal,decimal_count)
{
	var decimal_round = Math.round(decimal*Math.pow(10,decimal_count));
	return decimal_round/Math.pow(10,decimal_count);
}

function show_data(Byte_str)
{
	var data_count = parseInt(Byte_str);

	var data_str = '';
	//compute TB
	if(data_count >= 1024*1024*1024*1024)
	{
		var devide = data_count/(1024*1024*1024*1024);
		var reamian = data_count%(1024*1024*1024*1024);
		if(devide != 0)
		{
			data_str += show_decimal(devide,2)+' TB';
			return data_str;
		}
	}

	//compute GB
	if(data_count >= 1024*1024*1024)
	{
		var devide = data_count/(1024*1024*1024);
		var reamian = data_count%(1024*1024*1024);
		if(devide != 0)
		{
			data_str += show_decimal(devide,2)+' GB';
			return data_str;
		}
	}

	//compute MB
	if(data_count >= 1024*1024)
	{
		var devide = data_count/(1024*1024);
		var reamian = data_count%(1024*1024);
		if(devide != 0)
		{
			data_str += show_decimal(devide,2)+' MB';
			return data_str;
		}
	}

	//compute KB
	if(data_count >= 1024)
	{
		var devide = data_count/(1024);
		var reamian = data_count%(1024);
		if(devide != 0)
		{
			data_str += show_decimal(devide,2)+' KB';
			return data_str;
		}
	}

	//compute B
	var devide = data_count;
	data_str += devide+' Byte';
	return data_str;

}

function change_state()
{
	$.get('/goform/status_update?rd=' + Math.random(), function(r)
	{
		var status_info=r.split(";");
		show_net_state(status_info[1]);
		//document.getElementById('last_time').innerHTML = show_time(status_info[7]);
		show_connect_time(status_info[9],status_info[1]);
		document.getElementById('last_all').innerHTML = show_data(parseInt(status_info[10])+parseInt(status_info[11]));
		document.getElementById('tatal_all').innerHTML = show_data(parseInt(status_info[13])+parseInt(status_info[14]));
	    //gray(status_info[14]);

	});
}

function connect_net()
{
	var flag = true;
	var NetStat_roam = top.global_get_cookie("NetStat_roam");
	var roamstauts_roam = top.global_get_cookie("roamstauts_roam");
	var service_status_roam = top.global_get_cookie("service_status_roam");
	var roam_status_info = '0';
		if((NetStat_roam!='2') && ((roamstauts_roam == '1')&&('2' == service_status_roam)) && (roam_status_info == '1'))
		{
			if(!confirm(On_roam[lang_index]))
			{
				flag = false;
			}
		}
		if((NetStat_roam!='2') && ((roamstauts_roam == '1')&&('2' == service_status_roam)) && (roam_status_info!= '1'))
		{
			flag = false;
			alert(Roam_option_off[lang_index]);
		}
		if(flag == true)
		{
			top.global_set_cookie("disconnect_waiting_times","4");

			top.global_set_cookie("lastPage","1");
			if(Auto == '0')
			{
				top.global_set_cookie("always_on_to_manual","1");
				//$.get('/goform/SetNetworkSelectionMode?file=4&node=AutoManual&value=1&rd='+Math.random(),function(r){alert('h');});
				//$.ajax({url: '/goform/SetNetworkSelectionMode?file=4&node=AutoManual&value=1&rd='+Math.random(), async: false }).responseText; //??
				//alert("after");
			}
			document.netConnect.action = "/goform/ManualDisconnect?rd=" + Math.random();
		}
		else
		{
			return false;
		}
}
$(document).ready(function()
{
	$("#manual_mode_show").click(function()
	{
	});
	change_state();
	setInterval("change_state()",1000);
	setTimeout("change_state()",1);
});

function initpage()
{
	get_pin_status();
	initTranslation();
}
*/
