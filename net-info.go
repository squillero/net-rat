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
	"net"
	"time"
)

// ======================================================================
// Network info structures

type IpFlags int

const (
	LoopbackIP = 1 << iota
	TunnelIP
	LocalIP
	PublicIP
)

type IpInfo struct {
	RawIp     string
	Comment   string
	Source    string
	Flags     IpFlags
	Timestamp time.Time
}

func (ip IpInfo) Describe() string {
	cookedInfo := ip.RawIp
	if ip.Comment != "" {
		cookedInfo += " [" + ip.Comment + "]"
	}
	return cookedInfo
}

func (ip IpInfo) IsCool() bool {
	return ip.Comment != ""
}
func (ip IpInfo) IsValid() bool {
	return ip.RawIp != "" && time.Now().Sub(ip.Timestamp) < INFO_TIMEOUT
}

type NetInfo struct {
	ips map[string]IpInfo
}

func NewNetInfo() NetInfo {
	return NetInfo{ips: make(map[string]IpInfo)}
}

func (ni NetInfo) GetType(t IpFlags) []IpInfo {
	var r []IpInfo
	for _, v := range ni.ips {
		if v.Flags&t == t {
			r = append(r, v)
		}
	}
	return r
}

func checkKnownSubnets(ip IpInfo) string {
	knownNetworks := map[string]string{
		"130.192.0.0/16": "Politecnico di Torino",
	}
	ipt, _, _ := net.ParseCIDR(ip.RawIp + "/32")
	for k, v := range knownNetworks {
		_, subnet, _ := net.ParseCIDR(k)
		if subnet.Contains(ipt) {
			return v
		}
	}
	return ""
}

func (ni NetInfo) add(ip IpInfo) {
	if !ip.IsValid() {
		slog.Debug("Invalid IP", "ip", ip)
		return
	}
	if val, ok := ni.ips[ip.RawIp]; ok {
		if val.Flags == PublicIP || ip.Flags == PublicIP {
			ip.Flags = PublicIP
		}
		if len(val.Comment) > len(ip.Comment) {
			ip.Comment = val.Comment
		}
	}
	ni.ips[ip.RawIp] = ip
}

func (ni NetInfo) AnyCool(t IpFlags) bool {
	for _, v := range ni.ips {
		if v.Flags&t == t && v.IsCool() {
			return true
		}
	}
	return false
}

func (ni NetInfo) Any(t IpFlags) bool {
	for _, v := range ni.ips {
		if v.Flags&t == t {
			return true
		}
	}
	return false
}

// ======================================================================
// Get network info

func getNetInfo() NetInfo {
	result := NewNetInfo()

	// ... but first, cache
	if cache, err := CacheLoad(); err == nil {
		result = cache
		slog.Debug("Using cached data", "result", result)
	}

	// Create a context with a timeout of 1 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ipChan := make(chan IpInfo, 5)

	// Local IP info providers
	go getLocalIpIFACE(ipChan)
	go getLocalIpUDP(ipChan)
	//go getGateway(ipChan)

	// IsCool IP info providers
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
	for !timedOut && (!result.Any(LocalIP) || !result.AnyCool(PublicIP)) {
		select {
		case <-ctxTimeout.Done():
			slog.Debug("getNetInfo timeout!\n")
			timedOut = true
		case ip = <-ipChan:
			slog.Debug("Got IP info", "ip", ip.Describe(), "source", ip.Source)
			if info := checkKnownSubnets(ip); info != "" {
				ip.Comment = info
			}
			result.add(ip)
		}
	}
	CacheSave(result)
	return result
}
