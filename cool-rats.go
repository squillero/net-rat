//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

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
		slog.Error("getAirVPN", "error", err)
	}

	// check if "geo_additional" is present in cooked

	var geo string
	if geo_additional, ok := cooked["geo_additional"]; ok {
		geo = (geo_additional.(map[string]interface{}))["region_name"].(string) + ", " + (cooked["geo_additional"].(map[string]interface{}))["country_name"].(string)
	} else {
		geo = (cooked["geo"].(map[string]interface{}))["name"].(string)
	}
	info := IpInfo{
		RawIp:     cooked["ip"].(string),
		Comment:   geo,
		Source:    url,
		Flags:     PublicIP | CoolIP,
		Timestamp: time.Now(),
	}
	out <- info
}

func getIpGeoInfo(out chan IpInfo, url, ip, geoInfo, geoInfo2 string) {
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
		slog.Error("getIpGeoInfo", "error", err)
	}

	geo := cooked[geoInfo].(string)
	if geoInfo2 != "" {
		geo += ", " + cooked[geoInfo2].(string)
	}

	info := IpInfo{
		RawIp:     cooked[ip].(string),
		Comment:   geo,
		Source:    url,
		Flags:     PublicIP | CoolIP,
		Timestamp: time.Now(),
	}
	out <- info
}
