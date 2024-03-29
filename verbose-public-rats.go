//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   Yet another little Go experiment
//
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// This code is being released for educational and academic purposes.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PublicIpInfoVerbose struct {
	Ip        string
	Geo       string
	Source    string
	Timestamp time.Time
	reliable  bool
}

func (ni PublicIpInfoVerbose) String() string {
	ip := ni.Ip
	if ni.Geo != "" {
		ip += " (" + ni.Geo + ")"
	}
	return ip
}

func (ni PublicIpInfoVerbose) GetType() IpType {
	if !ni.reliable {
		return IllegalIP
	} else {
		return PublicIP
	}
}

// Fetch public IP info from Mullvad VPN provider public info checker
func getMullvad(out chan IpInfo) {
	const url string = "https://am.i.mullvad.net/json"

	result, err := http.Get(url)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}

	var cooked map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cooked); err != nil {
		// print out if error is not nil
		fmt.Println(err)
	}

	info := PublicIpInfoVerbose{
		Ip:     cooked["ip"].(string),
		Geo:    cooked["city"].(string) + "//" + cooked["country"].(string),
		Source: url,
	}
	time.Sleep(500 * time.Millisecond)
	out <- info
}

// Fetch public IP info from AirVPN VPN provider public info checker
func getAirVPN(out chan IpInfo) {
	const url string = "https://airvpn.org/api/whatismyip/"

	result, err := http.Get(url)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}

	var cooked map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cooked); err != nil {
		// print out if error is not nil
		fmt.Println(err)
	}

	// check if "geo_additional" is present in cooked

	var geo string
	if geo_additional, ok := cooked["geo_additional"]; ok {
		geo = (geo_additional.(map[string]interface{}))["region_name"].(string) + "/" + (cooked["geo_additional"].(map[string]interface{}))["country_name"].(string)
	} else {
		geo = (cooked["geo"].(map[string]interface{}))["name"].(string)
	}
	info := PublicIpInfoVerbose{
		Ip:     cooked["ip"].(string),
		Geo:    geo,
		Source: url,
	}
	time.Sleep(500 * time.Millisecond)
	out <- info
}
