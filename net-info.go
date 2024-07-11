//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"context"
	"log/slog"
	"time"
)

type IpFlags int

const (
	LocalIP  = 0x0
	PublicIP = 0x1
	IpType   = 0x01

	CoolIP = 0x10

	IPv6 IpFlags = 0x100
)

type IpInfo struct {
	RawIp     string
	CookedIp  string
	Source    string
	Flags     IpFlags
	Timestamp time.Time
}

func (ip IpInfo) String() string {
	return ip.CookedIp
}

func (ip IpInfo) Cool() bool {
	return ip.Flags&CoolIP == CoolIP
}
func (ip IpInfo) Valid() bool {
	return ip.RawIp != ""
}

func getNetInfo() [2]IpInfo {
	var result = [2]IpInfo{IpInfo{}, IpInfo{}}

	// ... but first, cache
	if cache, err := CacheLoad(); err == nil {
		for i := 0; i < 2; i++ {
			if time.Now().Sub(cache[i].Timestamp) < CACHE_TIMEOUT {
				result[cache[i].Flags&IpType] = cache[i]
			}
		}
		slog.Debug("Using cached data", "result", result)
	}

	// Create a context with a timeout of 1 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ipChan := make(chan IpInfo, 1)

	// Local IP info providers
	go getLocalIpIFACE(ipChan)
	go getLocalIpUDP(ipChan)
	//go getGateway(ipChan)

	// Cool IP info providers
	go getAirVPN(ipChan)
	go getIpGeoInfo(ipChan, "https://freeipapi.com/api/json", "ipAddress", "cityName", "countryCode")
	go getIpGeoInfo(ipChan, "https://am.i.mullvad.net/json", "ip", "city", "country")

	// Standard IP info providers
	go fetchRaw(ipChan, "http://api4.ipify.org/")
	go fetchRaw(ipChan, "https://checkip.amazonaws.com/")
	go fetchRaw(ipChan, "https://icanhazip.com/")
	go fetchRaw(ipChan, "http://ifconfig.me/ip")
	go fetchRaw(ipChan, "http://ipecho.net/plain")
	//go fetchJson(ipChan, "http://ipinfo.io", "ip")
	go fetchJson(ipChan, "http://ipv4.iplocation.net", "ip")

	var ip IpInfo
	timedOut := false
	for !timedOut && !(result[0].Cool() && result[1].Cool()) {
		select {
		case <-ctxTimeout.Done():
			slog.Debug("getNetInfo timeout!\n")
			timedOut = true
		case ip = <-ipChan:
			if result[ip.Flags&IpType].Valid() && result[ip.Flags&IpType].RawIp != ip.RawIp {
				slog.Debug("IP Mismatch!", "old", result[ip.Flags&IpType].RawIp, "new", ip.RawIp)
			}
			if !result[ip.Flags&IpType].Valid() || ip.Cool() {
				result[ip.Flags&IpType] = ip
			}
			slog.Debug("Got IP info", "ip", ip)
		}
	}
	CacheSave(result)
	return result
}
